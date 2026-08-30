package app

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"ava/pkg/logger"

	"ava/hub/internal/api"
	"ava/hub/internal/config"
	"ava/hub/internal/inventory"
	"ava/hub/internal/state"
	"ava/pkg/mqtt"
	"ava/pkg/wire"
)

const departTimeout = 2 * time.Second

type App struct {
	cfg    *config.Config
	client *api.Client
	// stateMu guards state, which the token loop replaces and the broker client
	// reads from its own goroutine every time it dials.
	stateMu  sync.RWMutex
	state    *state.State
	devices  *registry
	lastSeen []inventory.Entry
	mqtt     *mqtt.Client
	topics   wire.Topics
}

// brokerCredentials is what the hub currently believes its broker login to be.
//
// Handed to the client as a function rather than a value so a rotation picked
// up by the token loop reaches the next connection attempt. The API issues a
// fresh broker password every time it refreshes the hub's token, and for as
// long as this was read once at startup the hub spent the rest of its life
// presenting the password it was born with: online over HTTP, refused by the
// broker, and silent about both.
func (a *App) brokerCredentials() (username, password string) {
	a.stateMu.RLock()
	defer a.stateMu.RUnlock()

	if a.state == nil {
		return "", ""
	}

	return a.state.BrokerUsername, a.state.BrokerPassword
}

func Bootstrap(_ context.Context) (*App, error) {
	cfg := config.GetConfig()

	logger.Init(cfg.Env, cfg.LogLevel)

	loaded, err := state.Load(cfg.StateFile)
	if err != nil {
		return nil, err
	}

	logger.Info("HUB_STARTED",
		logger.String("hub_name", cfg.HubName),
		logger.String("api_base_url", cfg.APIBaseURL),
		logger.String("state_file", cfg.StateFile),
	)

	return &App{
		cfg:     cfg,
		client:  api.NewClient(cfg.APIBaseURL),
		state:   loaded,
		devices: newRegistry(),
	}, nil
}

func (a *App) Close() {
	if a == nil {
		return
	}

	a.depart()
	a.mqtt.Close()
	logger.Sync()
}

func (a *App) depart() {
	if a.mqtt == nil || a.topics.Status == "" {
		return
	}

	offline, err := json.Marshal(wire.Presence{Online: false, HubID: a.state.HubID})
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), departTimeout)
	defer cancel()

	if err := a.mqtt.Publish(ctx, a.topics.Status, offline, true); err != nil {
		logger.Warn("PRESENCE_PUBLISH_FAILED", logger.String("error", err.Error()))
	}
}

func (a *App) Run(ctx context.Context) error {
	tokens, err := a.authorize(ctx)
	if err != nil {
		return err
	}

	a.client.SetToken(tokens.AccessToken)

	logger.Info("HUB_PAIRED",
		logger.String("hub_id", tokens.Hub.ID),
		logger.String("hub_name", tokens.Hub.Name),
		logger.String("tenant", tokens.Tenant.Slug),
	)

	a.lastSeen = recall(a.state.Devices, a.cfg.DiscoveryTimeout)

	if _, err := a.startCommands(ctx, tokens); err != nil {
		logger.Warn("MQTT_UNAVAILABLE", logger.Err(err))
	}

	go a.syncLoop(ctx)

	return a.heartbeatLoop(ctx, tokens)
}

func (a *App) authorize(ctx context.Context) (*api.HubTokens, error) {
	if a.state.IsPaired() {
		tokens, err := a.client.RefreshToken(ctx, a.state.RefreshToken)
		if err == nil {
			return tokens, a.persist(tokens)
		}

		code := api.CodeOf(err)
		if code != api.CodeInvalidRefreshToken && code != api.CodeHubRevoked {
			return nil, err
		}

		logger.Warn("HUB_PAIRING_REJECTED", logger.String("code", code))

		if err := state.Clear(a.cfg.StateFile); err != nil {
			return nil, err
		}

		a.state = &state.State{}
	}

	return a.pair(ctx)
}

func (a *App) pair(ctx context.Context) (*api.HubTokens, error) {
	code, err := a.client.RequestActivationCode(ctx, a.cfg.HubName)
	if err != nil {
		return nil, err
	}

	logger.Info("HUB_ACTIVATION_REQUIRED",
		logger.String("user_code", code.UserCode),
		logger.String("verification_uri", code.VerificationURI),
		logger.Int64("expires_in", code.ExpiresIn),
	)

	interval := time.Duration(code.Interval) * time.Second
	if interval <= 0 {
		interval = 5 * time.Second
	}

	deadline := time.Now().Add(time.Duration(code.ExpiresIn) * time.Second)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(interval):
		}

		if time.Now().After(deadline) {
			return nil, errors.New("activation code expired before it was approved")
		}

		tokens, err := a.client.PollToken(ctx, code.DeviceCode)
		if err == nil {
			return tokens, a.persist(tokens)
		}

		switch api.CodeOf(err) {
		case api.CodeAuthorizationPending:
			continue
		case api.CodeSlowDown:
			interval += time.Second

			continue
		case api.CodeAccessDenied:
			return nil, errors.New("activation was denied")
		case api.CodeExpiredToken:
			return nil, errors.New("activation code expired before it was approved")
		default:
			return nil, err
		}
	}
}

func (a *App) persist(tokens *api.HubTokens) error {
	next := &state.State{
		HubID:        tokens.Hub.ID,
		HubName:      tokens.Hub.Name,
		TenantSlug:   tokens.Tenant.Slug,
		RefreshToken: tokens.RefreshToken,
	}

	a.stateMu.Lock()

	if a.state != nil {
		next.BrokerUsername = a.state.BrokerUsername
		next.BrokerPassword = a.state.BrokerPassword
		next.Devices = a.state.Devices
	}

	rotated := false

	if tokens.Broker != nil {
		rotated = tokens.Broker.Password != next.BrokerPassword
		next.BrokerUsername = tokens.Broker.Username
		next.BrokerPassword = tokens.Broker.Password
	}

	a.state = next
	saved := state.Save(a.cfg.StateFile, a.state)

	a.stateMu.Unlock()

	/* Worth a line of its own. A rotation is the moment the broker stops
	   accepting the connection the hub is holding, so it is the first thing to
	   look for when commands stop arriving. */
	if rotated {
		logger.Info("BROKER_CREDENTIALS_ROTATED", logger.String("username", next.BrokerUsername))
	}

	return saved
}

func (a *App) heartbeatLoop(ctx context.Context, tokens *api.HubTokens) error {
	renewAt := time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second / 2)

	ticker := time.NewTicker(a.cfg.HeartbeatInterval)
	defer ticker.Stop()

	for {
		connected := a.mqtt != nil && a.mqtt.Connected()

		if err := a.client.Heartbeat(ctx, connected); err != nil {
			logger.Warn("HEARTBEAT_FAILED", logger.String("error", err.Error()))
		}

		/* Said out loud on every beat, because the hub is perfectly capable of
		   looking healthy while being deaf. Heartbeats and device syncs are
		   HTTP and keep succeeding; commands arrive over the broker and do not.
		   Nothing else in the system notices the difference — the library is
		   silent about a refused reconnect, and the API's publish succeeds
		   whether or not anyone is subscribed. */
		if a.mqtt != nil && !connected {
			logger.Warn("COMMAND_CHANNEL_DOWN",
				logger.String("broker", a.cfg.MQTTBrokerURL),
				logger.String("detail", "commands cannot reach this hub until it reconnects"),
			)
		}

		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}

		if time.Now().After(renewAt) {
			refreshed, err := a.client.RefreshToken(ctx, a.state.RefreshToken)
			if err != nil {
				logger.Warn("TOKEN_REFRESH_FAILED", logger.String("error", err.Error()))

				continue
			}

			a.client.SetToken(refreshed.AccessToken)
			renewAt = time.Now().Add(time.Duration(refreshed.ExpiresIn) * time.Second / 2)

			if err := a.persist(refreshed); err != nil {
				return err
			}

			logger.Info("TOKEN_REFRESHED")
		}
	}
}
