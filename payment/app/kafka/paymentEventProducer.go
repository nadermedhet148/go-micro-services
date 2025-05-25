package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

type PaymentEvent struct {
	WALLET_ID int     `json:"wallet_id"`
	AMOUNT    float64 `json:"amount"`
	STATUS    string  `json:"status"`
	TYPE      string  `json:"type"` // "debit" or "credit"
}

type PaymentEventProducer struct {
	producer sarama.SyncProducer
}

func NewPaymentEventProducer() (PaymentEventProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	brokersUrl := []string{"localhost:29092"}

	producer, err := sarama.NewSyncProducer(brokersUrl, config)
	if err != nil {
		return PaymentEventProducer{}, err
	}

	return PaymentEventProducer{
		producer: producer,
	}, nil
}

func (ep *PaymentEventProducer) Close() error {
	return ep.producer.Close()
}

func (ep *PaymentEventProducer) PushToTopicWithPartition(message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic:     "payment-events",
		Value:     sarama.StringEncoder(message),
		Partition: -1, // -1 means the partition will be chosen by the producer
	}

	partition, offset, err := ep.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	fmt.Printf("Message is stored in partition(%d)/offset(%d)\n", partition, offset)
	return nil
}
