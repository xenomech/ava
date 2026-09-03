package broker

import (
	"context"
	"errors"

	"ava/pkg/mqtt"

	"github.com/google/uuid"
)

var ErrNotConnected = errors.New("broker: not connected")

type Config struct {
	URL      string
	ClientID string
	Username string
	Password string
	// CAFile trusts a private certificate authority; empty verifies against the system trust store.
	CAFile string
	// AllowInsecure permits plaintext to a public broker, which otherwise refuses to connect at all.
	AllowInsecure bool
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
		BrokerURL:     cfg.URL,
		ClientID:      clientID,
		Username:      cfg.Username,
		Password:      cfg.Password,
		CAFile:        cfg.CAFile,
		AllowInsecure: cfg.AllowInsecure,
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
