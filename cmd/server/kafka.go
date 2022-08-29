package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	config "gitlab.ozon.dev/pircuser61/catalog/config"
	counters "gitlab.ozon.dev/pircuser61/catalog/internal/counters"
	logger "gitlab.ozon.dev/pircuser61/catalog/internal/logger"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
	pkgCatalogConsumer "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/transport/grpc_kafka"
	"go.uber.org/zap"
)

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
	defer logger.Sync()
	ctx := context.Background()
	for {
		select {
		case <-session.Context().Done():
			logger.Debug("Kafka context channel closed")
			time.Sleep(time.Second * 10)
			return nil
		case msg, ok := <-claim.Messages():
			if !ok {
				logger.Debug("Kafka message channel closed ")
				return nil
			}
			logger.Debug("Kaffka new message",
				zap.Int32("partition", msg.Partition),
				zap.String("topic", msg.Topic),
				zap.ByteString("value", msg.Value))
			session.MarkMessage(msg, "")
			counters.Request()
			err := c.catalog.Handle(ctx, msg)
			if err != nil {
				counters.Error()
				logger.Error("Kafka response error", zap.Error(err))
				(*c.producer).Input() <- &sarama.ProducerMessage{
					Topic: config.Topic_error,
					Key:   sarama.StringEncoder(fmt.Sprint(time.Now())),
					Value: sarama.ByteEncoder([]byte(err.Error())),
				}
			} else {
				counters.Success()
			}
		}
	}
}

func runKafkaConsumer(ctx context.Context, core *storePkg.Core) {
	var topics = []string{config.Topic_create, config.Topic_update, config.Topic_delete}
	defer logger.Sync()

	cfg := sarama.NewConfig()
	asyncProducer, err := sarama.NewAsyncProducer(config.KafkaBrokers, cfg)
	if err != nil {
		return
	}

	kafkaCfg := sarama.NewConfig()
	kafkaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewConsumerGroup(config.KafkaBrokers, "catalog", kafkaCfg)
	if err != nil {
		logger.Panic("Kafka can't craete consumer: " + err.Error())
	}

	consumer := &Consumer{
		catalog:  pkgCatalogConsumer.New(core),
		producer: &asyncProducer,
	}
	for {
		logger.Debug("Consumer starts...")
		err := client.Consume(ctx, topics, consumer)
		if err != nil {
			logger.Error("Kaffka error on consume", zap.Error(err))
			time.Sleep(time.Second * 10)
		}
	}
}
