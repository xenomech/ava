package broker

import (
	"context"

	"ava/pkg/mqtt"

	"github.com/google/uuid"
)

type Config struct {
	URL      string
	ClientID string
}

type Broker struct {
	client *mqtt.Client
}

func Connect(ctx context.Context, cfg Config) (*Broker, error) {
	clientID := cfg.ClientID
	if clientID == "" {
		clientID = "ava-api-" + uuid.NewString()
	}

	client, err := mqtt.Connect(ctx, &mqtt.Options{
		BrokerURL: cfg.URL,
		ClientID:  clientID,
	})
	if err != nil {
		return nil, err
	}

	return &Broker{client: client}, nil
}

func (b *Broker) Publish(ctx context.Context, topic string, payload []byte, retained bool) error {
	if b == nil {
		return mqtt.ErrNotConnected
	}

	return b.client.Publish(ctx, topic, payload, retained)
}

func (b *Broker) Close() {
	if b == nil {
		return
	}

	b.client.Close()
}
