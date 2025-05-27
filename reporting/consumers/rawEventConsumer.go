package consumers

// SIGUSR1 toggle the pause/resume consumption
import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	dbConfig "example.com/reporting/config"

	"github.com/IBM/sarama"
)

// Sarama configuration options
var (
	brokers = "localhost:29092"
	group   = "raw-consumer-group"
	topics  = "payment-events"
	oldest  = true
	verbose = false
)

// -- `default`.payment_events definition

// CREATE TABLE default.payment_events
// (

//     `wallet_id` Int32,

//     `amount` Float64,

//     `status` String,

//     `type` String
// )
// ENGINE = MergeTree
// ORDER BY wallet_id
// SETTINGS index_granularity = 8192;

type PaymentEvent struct {
	WALLET_ID int     `json:"wallet_id"`
	AMOUNT    float64 `json:"amount"`
	STATUS    string  `json:"status"`
	TYPE      string  `json:"type"` // "debit" or "credit"
}

func RunGroup() {
	keepRunning := true
	log.Println("Starting a new Sarama consumer")

	if verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	/**
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */
	config := sarama.NewConfig()

	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	/**
	 * Setup a new Sarama consumer group
	 */

	conn, err := dbConfig.DbConnect()
	if err != nil {
		log.Panicf("Error connecting to ClickHouse: %v", err)
	}
	consumer := Consumer{
		ready: make(chan bool),
		conn:  conn, // Assign the ClickHouse connection to the consumer
	}

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(strings.Split(brokers, ","), group, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, strings.Split(topics, ","), &consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		case <-sigusr1:
			toggleConsumptionFlow(client, &consumptionIsPaused)
		}
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		log.Println("Resuming consumption")
	} else {
		client.PauseAll()
		log.Println("Pausing consumption")
	}

	*isPaused = !*isPaused
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool
	conn  driver.Conn // ClickHouse connection, if needed
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed")
				return nil
			}
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)

			event := PaymentEvent{}
			// Here you can unmarshal the message into your PaymentEvent struct
			// For example, if the message is in JSON format:
			err := json.Unmarshal(message.Value, &event)
			if err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				return err
			}
			// Here you can process the message, e.g., insert it into ClickHouse
			// Example: Insert into ClickHouse (assuming you have a connection)
			err = consumer.conn.Exec(context.Background(), `
				INSERT INTO default.payment_events
				(wallet_id, amount, status, type)
				VALUES (?, ?, ?, ?)`,
				event.WALLET_ID, event.AMOUNT, event.STATUS, event.TYPE,
			)
			if err != nil {
				log.Printf("Error inserting into ClickHouse: %v", err)
				return err
			}
			// If you have a specific ClickHouse connection, use it to execute the insert
			// For example, if you have a connection `conn`:
			// conn.Exec(context.Background(), `

			session.MarkMessage(message, "")
		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/IBM/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}
