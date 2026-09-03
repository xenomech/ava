package mqtt

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"ava/pkg/logger"

	paho "github.com/eclipse/paho.mqtt.golang"
)

const (
	connectTimeout = 10 * time.Second
	actionTimeout  = 5 * time.Second

	QoSAtLeastOnce = 1
)

var (
	ErrNoBroker     = errors.New("mqtt: broker url is required")
	ErrNotConnected = errors.New("mqtt: not connected")
)

type Handler func(topic string, payload []byte)

type Options struct {
	BrokerURL string
	ClientID  string
	Username  string
	Password  string
	// Credentials, when set, is asked afresh on every connection attempt so a rotated password survives.
	Credentials func() (username, password string)
	WillTopic   string
	Will        []byte
	Durable     bool
	OnConnect   func(client *Client)
	// CAFile trusts a private certificate authority; empty verifies against the system trust store.
	CAFile string
	// AllowInsecure permits plaintext to a public broker, which otherwise refuses to connect at all.
	AllowInsecure bool
}

type Client struct {
	inner paho.Client
	mu    sync.Mutex
	subs  map[string]Handler
}

func Connect(ctx context.Context, opts *Options) (*Client, error) {
	if opts == nil || opts.BrokerURL == "" {
		return nil, ErrNoBroker
	}

	if err := guardTransport(opts.BrokerURL, opts.AllowInsecure); err != nil {
		return nil, err
	}

	tlsConfig, err := tlsConfigFor(opts.CAFile)
	if err != nil {
		return nil, err
	}

	client := &Client{subs: make(map[string]Handler)}

	config := paho.NewClientOptions().
		AddBroker(opts.BrokerURL).
		SetClientID(opts.ClientID).
		SetCleanSession(!opts.Durable).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(5 * time.Second).
		// Short, because a drop here is usually a password rotation and the reconnect restores commands.
		SetMaxReconnectInterval(15 * time.Second).
		SetKeepAlive(30 * time.Second).
		SetConnectionLostHandler(func(_ paho.Client, err error) {
			logger.Warn("MQTT_DISCONNECTED", logger.Err(err))
		}).
		SetOnConnectHandler(func(_ paho.Client) {
			logger.Info("MQTT_CONNECTED", logger.String("broker", opts.BrokerURL))

			client.resume()

			if opts.OnConnect != nil {
				opts.OnConnect(client)
			}
		})

	if tlsConfig != nil {
		config = config.SetTLSConfig(tlsConfig)
	}

	switch {
	case opts.Credentials != nil:
		config = config.SetCredentialsProvider(paho.CredentialsProvider(opts.Credentials))
	case opts.Username != "":
		config = config.SetUsername(opts.Username).SetPassword(opts.Password)
	}

	if opts.WillTopic != "" {
		config = config.SetBinaryWill(opts.WillTopic, opts.Will, QoSAtLeastOnce, true)
	}

	client.inner = paho.NewClient(config)

	if err := wait(ctx, client.inner.Connect(), connectTimeout); err != nil {
		return nil, fmt.Errorf("mqtt: connect to %s: %w", opts.BrokerURL, err)
	}

	return client, nil
}

// Subscribe records the topic as wanted even when the subscribe fails, so resume can replay it later.
func (c *Client) Subscribe(ctx context.Context, topic string, handler Handler) error {
	c.mu.Lock()
	c.subs[topic] = handler
	c.mu.Unlock()

	if !c.ready() {
		return ErrNotConnected
	}

	return c.subscribe(ctx, topic, handler)
}

func (c *Client) subscribe(ctx context.Context, topic string, handler Handler) error {
	token := c.inner.Subscribe(topic, QoSAtLeastOnce, func(_ paho.Client, message paho.Message) {
		handler(message.Topic(), message.Payload())
	})

	if err := wait(ctx, token, actionTimeout); err != nil {
		return fmt.Errorf("mqtt: subscribe to %s: %w", topic, err)
	}

	logger.Info("MQTT_SUBSCRIBED", logger.String("topic", topic))

	return nil
}

func (c *Client) resume() {
	c.mu.Lock()
	pending := make(map[string]Handler, len(c.subs))

	for topic, handler := range c.subs {
		pending[topic] = handler
	}
	c.mu.Unlock()

	for topic, handler := range pending {
		ctx, cancel := context.WithTimeout(context.Background(), actionTimeout)

		if err := c.subscribe(ctx, topic, handler); err != nil {
			logger.Warn("MQTT_RESUBSCRIBE_FAILED", logger.String("topic", topic), logger.Err(err))
		}

		cancel()
	}
}

func (c *Client) Publish(ctx context.Context, topic string, payload []byte, retained bool) error {
	if !c.ready() {
		return ErrNotConnected
	}

	if err := wait(ctx, c.inner.Publish(topic, QoSAtLeastOnce, retained, payload), actionTimeout); err != nil {
		return fmt.Errorf("mqtt: publish to %s: %w", topic, err)
	}

	return nil
}

// Connected reports whether the client holds a broker session, since a refused reconnect is silent.
func (c *Client) Connected() bool {
	return c != nil && c.inner != nil && c.inner.IsConnected()
}

func (c *Client) Close() {
	if c != nil && c.inner != nil {
		c.inner.Disconnect(250)
	}
}

func (c *Client) ready() bool {
	return c != nil && c.inner != nil && c.inner.IsConnected()
}

func wait(ctx context.Context, token paho.Token, timeout time.Duration) error {
	select {
	case <-token.Done():
		return token.Error()
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(timeout):
		return errors.New("timed out")
	}
}
