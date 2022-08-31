package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
	goodFxtr "gitlab.ozon.dev/pircuser61/catalog/tests/fixtures"
)

func getSomeGood(code uint64) *models.Good {
	return goodFxtr.Good().Code(code).Name("name").UnitOfMeasure("uom").Country("country").P()
}

func TestHandler(t *testing.T) {
	ctx := context.Background()

	t.Run("topic found", func(t *testing.T) {
		// arrange
		f := catalogSetUp(ctx, t)
		handler := func(context.Context, *CatalogConsumer, *[]byte) error { return nil }
		f.kafkaConsumer.register("some_func", handler)
		msg := &sarama.ConsumerMessage{Topic: "some_func"}
		//act
		err := f.kafkaConsumer.Handle(ctx, msg)

		//assert
		assert.NoError(t, err)
	})
	t.Run("topic not found", func(t *testing.T) {
		// arrange
		f := catalogSetUp(ctx, t)
		handler := func(context.Context, *CatalogConsumer, *[]byte) error { return nil }

		f.kafkaConsumer.register("some_func", handler)
		msg := &sarama.ConsumerMessage{Topic: "another func"}

		//act
		err := f.kafkaConsumer.Handle(ctx, msg)

		//assert
		assert.Errorf(t, err, "Invalid topic name: %s", msg.Topic)
	})
}

func TestNewGoodService(t *testing.T) {
	ctx := context.Background()
	getReqMsg := func(good *models.Good) *sarama.ConsumerMessage {
		js, _ := json.Marshal(good)
		return &sarama.ConsumerMessage{Topic: "good_create", Value: js}
	}

	t.Run("good add success", func(t *testing.T) {
		// arrange
		good := getSomeGood(0)
		reqMsg := getReqMsg(good)

		f := catalogSetUp(ctx, t)
		f.goodRepo.EXPECT().Add(ctx, good).Return(nil)

		// act
		err := f.kafkaConsumer.Handle(ctx, reqMsg)

		// assert
		require.NoError(t, err)
		//assert.Equal(t, &pb.GoodCreateResponse{}, resp)
	})

	t.Run("create error ", func(t *testing.T) {

		cases := []struct {
			name  string
			errIn error
		}{{"keys not found", models.ErrValidation},
			{"timeout", storePkg.ErrTimeout},
			{"internal error", errors.New("error text")}}

		for _, tc := range cases {
			itc := tc
			t.Run(itc.name, func(t *testing.T) {
				// arrange
				good := getSomeGood(0)
				reqMsg := getReqMsg(good)
				f := catalogSetUp(ctx, t)
				f.goodRepo.EXPECT().Add(ctx, good).Return(itc.errIn)

				// act
				err := f.kafkaConsumer.Handle(ctx, reqMsg)

				// assert
				assert.ErrorIs(t, err, itc.errIn)
			})
		}
	})
}

func TestUpdateGoodService(t *testing.T) {
	ctx := context.Background()
	getReqMsg := func(good *models.Good) *sarama.ConsumerMessage {
		js, _ := json.Marshal(good)
		return &sarama.ConsumerMessage{Topic: "good_update", Value: js}
	}

	t.Run("good update success", func(t *testing.T) {
		// arrange
		good := getSomeGood(1)
		reqMsg := getReqMsg(good)

		f := catalogSetUp(ctx, t)
		f.goodRepo.EXPECT().Update(ctx, good).Return(nil)
		// act
		err := f.kafkaConsumer.Handle(ctx, reqMsg)

		// assert
		require.NoError(t, err)
	})

	t.Run("update error ", func(t *testing.T) {
		cases := []struct {
			name  string
			errIn error
		}{{"not found", storePkg.ErrNotExists},
			{"keys not found", models.ErrValidation},
			{"timeout", storePkg.ErrTimeout},
			{"internal error", errors.New("error text")}}

		for _, tc := range cases {
			itc := tc
			t.Run(itc.name, func(t *testing.T) {
				// arrange
				good := getSomeGood(0)
				reqMsg := getReqMsg(good)
				f := catalogSetUp(ctx, t)
				f.goodRepo.EXPECT().Update(ctx, good).Return(itc.errIn)

				// act
				err := f.kafkaConsumer.Handle(ctx, reqMsg)

				// assert
				assert.ErrorIs(t, err, itc.errIn)
			})
		}
	})
}

func TestGoodDeleteService(t *testing.T) {
	ctx := context.Background()
	code := uint64(22)

	getReqMsg := func(code uint64) *sarama.ConsumerMessage {
		js, _ := json.Marshal(code)
		return &sarama.ConsumerMessage{Topic: "good_delete", Value: js}
	}

	t.Run("good delete success", func(t *testing.T) {
		// arrange

		reqMsg := getReqMsg(code)
		f := catalogSetUp(ctx, t)
		f.goodRepo.EXPECT().Delete(ctx, code).Return(nil)

		// act
		err := f.kafkaConsumer.Handle(ctx, reqMsg)

		// assert
		require.NoError(t, err)
	})

	t.Run("good delete error ", func(t *testing.T) {
		cases := []struct {
			name  string
			errIn error
		}{{"not found", storePkg.ErrNotExists},
			{"timeout", storePkg.ErrTimeout},
			{"internal error", errors.New("error text")}}

		for _, tc := range cases {
			itc := tc
			t.Run(itc.name, func(t *testing.T) {
				// arrange
				reqMsg := getReqMsg(code)
				f := catalogSetUp(ctx, t)
				f.goodRepo.EXPECT().Delete(ctx, code).Return(itc.errIn)

				// act
				err := f.kafkaConsumer.Handle(ctx, reqMsg)

				// assert
				assert.ErrorIs(t, err, itc.errIn)
			})
		}
	})
}
