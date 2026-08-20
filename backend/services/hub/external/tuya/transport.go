package tuya

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	Port           = 6668
	defaultTimeout = 5 * time.Second
	readBuffer     = 4096
)

type transport interface {
	Do(ctx context.Context, addr string, request []byte) ([]byte, error)
}

type tcpTransport struct {
	timeout time.Duration
}

func newTCPTransport(timeout time.Duration) *tcpTransport {
	if timeout <= 0 {
		timeout = defaultTimeout
	}

	return &tcpTransport{timeout: timeout}
}

func (t *tcpTransport) Do(ctx context.Context, addr string, request []byte) ([]byte, error) {
	dialer := net.Dialer{Timeout: t.timeout}

	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("tuya: dial %s: %w", addr, err)
	}

	defer conn.Close()

	deadline := time.Now().Add(t.timeout)
	if got, ok := ctx.Deadline(); ok && got.Before(deadline) {
		deadline = got
	}

	if err := conn.SetDeadline(deadline); err != nil {
		return nil, fmt.Errorf("tuya: set deadline: %w", err)
	}

	if _, err := conn.Write(request); err != nil {
		return nil, fmt.Errorf("tuya: write: %w", err)
	}

	return readFrame(conn)
}

func readFrame(conn io.Reader) ([]byte, error) {
	header := make([]byte, headerLen)

	if _, err := io.ReadFull(conn, header); err != nil {
		return nil, fmt.Errorf("tuya: read header: %w", err)
	}

	length := int(binary.BigEndian.Uint32(header[12:16]))
	if length > maxPayloadLen {
		return nil, fmt.Errorf("%w: %d", errPayloadLimit, length)
	}

	body := make([]byte, length)

	if _, err := io.ReadFull(conn, body); err != nil {
		return nil, fmt.Errorf("tuya: read body: %w", err)
	}

	return append(header, body...), nil
}

func address(ip string) string {
	return net.JoinHostPort(ip, fmt.Sprint(Port))
}
