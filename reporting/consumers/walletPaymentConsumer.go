package consumers

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/gmbyapa/kstream/v2/kafka"
	"github.com/gmbyapa/kstream/v2/streams"
	"github.com/gmbyapa/kstream/v2/streams/encoding"
	"github.com/tryfix/log"
)

const TopicTextLines = `payment-events`
const TopicWordCount = `wallet-payments-count`

func RunStream() {
	flag.Parse()

	config := streams.NewStreamBuilderConfig()
	config.BootstrapServers = []string{`localhost:29092`}
	config.ApplicationId = `wallet-payments-count2`
	config.Consumer.Offsets.Initial = kafka.OffsetEarliest
	config.Consumer.GroupId = `wallet-payments-count-group2`
	config.Store.Changelog.ReplicaCount = 1
	config.InternalTopicsDefaultReplicaCount = 1

	builder := streams.NewStreamBuilder(config)
	buildTopology(config.Logger, builder)

	topology, err := builder.Build()
	if err != nil {
		panic(err)
	}

	println("Topology - \n", topology.Describe())

	runner := builder.NewRunner()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	go func() {
		<-sigs
		if err := runner.Stop(); err != nil {
			println(err)
		}
	}()

	if err := runner.Run(topology); err != nil {
		panic(err)
	}
}

func buildTopology(logger log.Logger, builder *streams.StreamBuilder) {
	stream := builder.KStream(TopicTextLines, encoding.StringEncoder{}, encoding.StringEncoder{})
	stream.Each(func(ctx context.Context, key, value any) {
		logger.Debug(`Processing event: ` + value.(string))
	}).SelectKey(func(ctx context.Context, key, value any) (kOut any, err error) {
		// Use WALLET_ID as the key
		var evt PaymentEvent
		if err := json.Unmarshal([]byte(value.(string)), &evt); err != nil {
			return nil, err
		}
		// Return the wallet ID as the key for grouping
		return fmt.Sprintf("%v", evt.WALLET_ID), nil
	}).Aggregate(TopicWordCount,
		func(ctx context.Context, key, value, previous any) (newAgg any, err error) {
			var count int
			if previous != nil {
				count = previous.(int)
			}
			count++
			newAgg = count
			return
		}, streams.AggregateWithValEncoder(encoding.IntEncoder{})).ToStream().Each(func(ctx context.Context, key, value any) {
		println(fmt.Sprintf(`WalletID %s: %d payments`, key, value))
	}).To(TopicWordCount)
}
