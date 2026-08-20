package wiz

import (
	"context"
	"fmt"
	"net"
	"time"
)

const (
	Port           = 38899
	maxAttempts    = 4
	initialBackoff = 250 * time.Millisecond
	readBuffer     = 2048
)

type transport interface {
	Do(ctx context.Context, addr string, payload []byte) ([]byte, error)
}

type udpTransport struct {
	timeout time.Duration
}

func newUDPTransport(timeout time.Duration) *udpTransport {
	if timeout <= 0 {
		timeout = time.Second
	}

	return &udpTransport{timeout: timeout}
}

func (t *udpTransport) Do(ctx context.Context, addr string, payload []byte) ([]byte, error) {
	var lastErr error

	backoff := initialBackoff

	for attempt := range maxAttempts {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}

			backoff *= 2
		}

		reply, err := t.exchange(ctx, addr, payload)
		if err == nil {
			return reply, nil
		}

		lastErr = err

		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
	}

	return nil, fmt.Errorf("wiz: no reply from %s after %d attempts: %w", addr, maxAttempts, lastErr)
}

func (t *udpTransport) exchange(ctx context.Context, addr string, payload []byte) ([]byte, error) {
	dialer := net.Dialer{Timeout: t.timeout}

	conn, err := dialer.DialContext(ctx, "udp", addr)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	defer conn.Close()

	deadline := time.Now().Add(t.timeout)
	if got, ok := ctx.Deadline(); ok && got.Before(deadline) {
		deadline = got
	}

	if err := conn.SetDeadline(deadline); err != nil {
		return nil, fmt.Errorf("set deadline: %w", err)
	}

	if _, err := conn.Write(payload); err != nil {
		return nil, fmt.Errorf("write: %w", err)
	}

	buf := make([]byte, readBuffer)

	n, err := conn.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	return buf[:n], nil
}

func address(ip string) string {
	return net.JoinHostPort(ip, fmt.Sprint(Port))
}
