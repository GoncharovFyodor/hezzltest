package nats

import (
	"context"
	"errors"
	"fmt"
	"github.com/GoncharovFyodor/hezzltest/internal/config"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"log"
)

func New(ctx context.Context, cfg *config.Config) (*nats.Conn, jetstream.JetStream, jetstream.Stream) {
	nc, err := nats.Connect(fmt.Sprintf("nats://%s", cfg.Nats.Host))

	if err != nil {
		log.Fatalf("не удалось подключиться к nats: %v", err)
	}

	js, err := jetstream.New(nc)

	if err != nil {
		log.Fatalf("не удалось создать jetstream: %v", err)
	}

	stream, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     "goods",
		Subjects: []string{"goods.*"},
	})

	if !errors.Is(err, jetstream.ErrStreamNameAlreadyInUse) && err != nil {
		log.Fatalf("не удалось создать stream: %v", err)
	}

	return nc, js, stream
}
