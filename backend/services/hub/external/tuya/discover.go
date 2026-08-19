package tuya

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"ava/hub/internal/device"
)

const (
	DiscoveryPort          = 6666
	DiscoveryPortEncrypted = 6667

	discoverySecret = "yGAdlopoPVldABfn"
)

type Found struct {
	Info    device.Info
	Version string
	Product string
}

func (f *Found) Supported() bool {
	return f.Version == Version
}

type broadcast struct {
	IP         string `json:"ip"`
	GatewayID  string `json:"gwId"`
	ProductKey string `json:"productKey"`
	Version    string `json:"version"`
	Active     int    `json:"active"`
}

func discoveryKey() []byte {
	sum := md5.Sum([]byte(discoverySecret))

	return sum[:]
}

func Discover(ctx context.Context, timeout time.Duration) ([]Found, error) {
	if timeout <= 0 {
		timeout = 6 * time.Second
	}

	deadline := time.Now().Add(timeout)
	if got, ok := ctx.Deadline(); ok && got.Before(deadline) {
		deadline = got
	}

	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		seen    = make(map[string]struct{})
		found   = make([]Found, 0)
		lastErr error
	)

	for _, port := range []int{DiscoveryPort, DiscoveryPortEncrypted} {
		wg.Add(1)

		go func() {
			defer wg.Done()

			results, err := listen(port, deadline)
			if err != nil {
				mu.Lock()
				lastErr = err
				mu.Unlock()

				return
			}

			mu.Lock()
			defer mu.Unlock()

			for at := range results {
				if _, duplicate := seen[results[at].Info.ID]; duplicate {
					continue
				}

				seen[results[at].Info.ID] = struct{}{}
				found = append(found, results[at])
			}
		}()
	}

	wg.Wait()

	if len(found) == 0 && lastErr != nil {
		return nil, lastErr
	}

	return found, nil
}

func listen(port int, deadline time.Time) ([]Found, error) {
	conn, err := net.ListenPacket("udp4", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("tuya: listen on %d: %w", port, err)
	}

	defer conn.Close()

	if err := conn.SetDeadline(deadline); err != nil {
		return nil, fmt.Errorf("tuya: set deadline: %w", err)
	}

	found := make([]Found, 0)
	buf := make([]byte, readBuffer)

	for {
		n, _, err := conn.ReadFrom(buf)
		if err != nil {
			return found, nil
		}

		announcement, err := decodeBroadcast(buf[:n])
		if err != nil {
			continue
		}

		found = append(found, announcement)
	}
}

func decodeBroadcast(data []byte) (Found, error) {
	payload := data

	if frame, err := unpack(data, true); err == nil {
		payload = frame.payload
	}

	if !looksLikeJSON(payload) {
		decrypted, err := decrypt(discoveryKey(), payload)
		if err != nil {
			return Found{}, fmt.Errorf("tuya: decrypt broadcast: %w", err)
		}

		payload = decrypted
	}

	var message broadcast
	if err := json.Unmarshal(payload, &message); err != nil {
		return Found{}, fmt.Errorf("tuya: decode broadcast: %w", err)
	}

	if message.GatewayID == "" {
		return Found{}, errors.New("tuya: broadcast has no device id")
	}

	return Found{
		Info: device.Info{
			ID:       message.GatewayID,
			Vendor:   Vendor,
			IP:       message.IP,
			LastSeen: time.Now(),
		},
		Version: message.Version,
		Product: message.ProductKey,
	}, nil
}

func looksLikeJSON(payload []byte) bool {
	return len(payload) > 1 && payload[0] == '{' && payload[len(payload)-1] == '}'
}
