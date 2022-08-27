package kafka

import (
	"context"

	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"

	"gitlab.ozon.dev/pircuser61/catalog/config"
	log "gitlab.ozon.dev/pircuser61/catalog/internal/log"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
)

type msgHandler func(context.Context, *CatalogConsumer, *[]byte) error

type CatalogConsumer struct {
	route map[string]msgHandler
	store *storePkg.Core
}

func New(store *storePkg.Core) *CatalogConsumer {
	cc := CatalogConsumer{route: make(map[string]msgHandler), store: store}
	cc.register(config.Topic_create, goodCreate)
	cc.register(config.Topic_update, goodUpdate)
	cc.register(config.Topic_delete, goodDelete)
	return &cc
}

func (cc *CatalogConsumer) register(topic string, handler msgHandler) {
	cc.route[topic] = handler
}

func (cc *CatalogConsumer) Handle(ctx context.Context, msg *sarama.ConsumerMessage) error {
	handler, ok := cc.route[msg.Topic]
	if ok {
		return handler(ctx, cc, &msg.Value)
	} else {
		log.Warnf("Invalid topic name: %s", msg.Topic)
		return errors.Errorf("Invalid topic name: %s", msg.Topic)
	}
}

func goodCreate(ctx context.Context, cc *CatalogConsumer, in *[]byte) error {
	good := &models.Good{}
	if err := json.Unmarshal(*in, good); err != nil {
		return err
	}
	log.Msgf("Good create %v", good)
	return cc.store.Good.Add(ctx, good)
}

func goodUpdate(ctx context.Context, cc *CatalogConsumer, in *[]byte) error {
	good := &models.Good{}
	if err := json.Unmarshal(*in, good); err != nil {
		return err
	}
	log.Msgf("Good update %v", good)
	return cc.store.Good.Update(ctx, good)
}

func goodDelete(ctx context.Context, cc *CatalogConsumer, in *[]byte) error {
	var code uint64
	if err := json.Unmarshal(*in, &code); err != nil {
		return err
	}
	log.Msgf("Good delete %v", code)
	return cc.store.Good.Delete(ctx, code)
}
