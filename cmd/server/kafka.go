package main

import (
	"context"
	"expvar"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	config "gitlab.ozon.dev/pircuser61/catalog/config"
	log "gitlab.ozon.dev/pircuser61/catalog/internal/log"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
	pkgCatalogConsumer "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/transport/grpc_kafka"
)

var counterIn, counterSuccess, counterErr expvar.Int

type Consumer struct {
	catalog  *pkgCatalogConsumer.CatalogConsumer
	producer *sarama.AsyncProducer
}

func (c *Consumer) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx := context.Background()
	for {
		select {
		case <-session.Context().Done():
			log.Msg("Done")
			return nil
		case msg, ok := <-claim.Messages():
			if !ok {
				log.Msg("Data channel closed")
				return nil
			}
			log.Msgf("partition: %v\n, topic: %v\n data: %v", msg.Partition, msg.Topic, string(msg.Value))
			session.MarkMessage(msg, "")
			counterIn.Add(1)
			err := c.catalog.Handle(ctx, msg)
			if err != nil {
				counterErr.Add(1)
				log.ErrorMsg(err.Error())
				(*c.producer).Input() <- &sarama.ProducerMessage{
					Topic: config.Topic_error,
					Key:   sarama.StringEncoder(fmt.Sprint(time.Now())),
					Value: sarama.ByteEncoder([]byte(err.Error())),
				}
			} else {
				counterSuccess.Add(1)
			}
		}
	}
}

func runKafkaConsumer(ctx context.Context, core *storePkg.Core) {
	var topics = []string{config.Topic_create, config.Topic_update, config.Topic_delete}

	cfg := sarama.NewConfig()
	asyncProducer, err := sarama.NewAsyncProducer(config.KafkaBrokers, cfg)
	if err != nil {
		return
	}

	kafkaCfg := sarama.NewConfig()
	kafkaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewConsumerGroup(config.KafkaBrokers, "catalog", kafkaCfg)
	if err != nil {
		log.Msgf("ERR on create kafka client: %v", err)
		return
	}

	consumer := &Consumer{
		catalog:  pkgCatalogConsumer.New(core),
		producer: &asyncProducer,
	}
	for {
		log.Msg("Consumer starts...")
		err := client.Consume(ctx, topics, consumer)
		if err != nil {
			log.Msgf("ERR on consume: %v", err)
			time.Sleep(time.Second * 10)
		}
	}
}
