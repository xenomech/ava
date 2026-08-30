package wiz

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"ava/hub/internal/device"
	"ava/pkg/wire"
)

const (
	BroadcastAddr = "255.255.255.255"

	probeInterval = 500 * time.Millisecond
)

type Found struct {
	Info  device.Info
	State wire.State
}

func Discover(ctx context.Context, timeout time.Duration) ([]Found, error) {
	if timeout <= 0 {
		timeout = 3 * time.Second
	}

	conn, err := net.ListenPacket("udp4", ":0")
	if err != nil {
		return nil, fmt.Errorf("wiz: listen: %w", err)
	}

	defer conn.Close()

	deadline := time.Now().Add(timeout)
	if got, ok := ctx.Deadline(); ok && got.Before(deadline) {
		deadline = got
	}

	if err := conn.SetDeadline(deadline); err != nil {
		return nil, fmt.Errorf("wiz: set deadline: %w", err)
	}

	target, err := net.ResolveUDPAddr("udp4", address(BroadcastAddr))
	if err != nil {
		return nil, fmt.Errorf("wiz: resolve broadcast: %w", err)
	}

	probe, err := json.Marshal(request{Method: "getPilot", Params: struct{}{}})
	if err != nil {
		return nil, fmt.Errorf("wiz: encode probe: %w", err)
	}

	if _, err := conn.WriteTo(probe, target); err != nil {
		return nil, fmt.Errorf("wiz: broadcast: %w", err)
	}

	go reprobe(ctx, conn, target, probe, deadline)

	return collect(ctx, conn), nil
}

func reprobe(ctx context.Context, conn net.PacketConn, target net.Addr, probe []byte, deadline time.Time) {
	ticker := time.NewTicker(probeInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}

		if !time.Now().Before(deadline) {
			return
		}

		if _, err := conn.WriteTo(probe, target); err != nil {
			return
		}
	}
}

func collect(ctx context.Context, conn net.PacketConn) []Found {
	seen := make(map[string]struct{})
	found := make([]Found, 0)
	buf := make([]byte, readBuffer)

	for {
		if ctx.Err() != nil {
			return found
		}

		n, from, err := conn.ReadFrom(buf)
		if err != nil {
			return found
		}

		ip, _, splitErr := net.SplitHostPort(from.String())
		if splitErr != nil {
			continue
		}

		if _, duplicate := seen[ip]; duplicate {
			continue
		}

		var envelope response
		if err := json.Unmarshal(buf[:n], &envelope); err != nil || envelope.Error != nil {
			continue
		}

		var result pilotResult
		if err := json.Unmarshal(envelope.Result, &result); err != nil {
			continue
		}

		seen[ip] = struct{}{}

		found = append(found, Found{
			Info: device.Info{
				ID:       result.MAC,
				Vendor:   Vendor,
				IP:       ip,
				MAC:      result.MAC,
				LastSeen: time.Now(),
			},
			State: pilotState(&result),
		})
	}
}

func pilotState(result *pilotResult) wire.State {
	state := wire.State{wire.TraitPower: wire.Bool(result.State)}

	if result.Dimming > 0 {
		state[wire.TraitBrightness] = wire.Number(float64(result.Dimming))
	}

	if result.Temp > 0 {
		state[wire.TraitColorTemp] = wire.Number(float64(result.Temp))
	}

	return state
}
