package nats

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go/jetstream"
	log "github.com/sirupsen/logrus"
)

type Producer struct {
	log *log.Logger
	js  jetstream.JetStream
}

func NewProducer(js jetstream.JetStream) *Producer {
	return &Producer{js: js}
}

func (p *Producer) Publish(ctx context.Context, good any) error {
	bytes, err := json.Marshal(&good)

	if err != nil {
		p.log.Infof("ошибка при сериализации товара: %v", err)
		return err
	}

	_, err = p.js.Publish(ctx, "goods.test", bytes)

	if err != nil {
		p.log.Infof("ошибка при отправке товара в nats: %v", err)
	}

	return nil
}
