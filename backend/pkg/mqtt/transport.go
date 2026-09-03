package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
)

// ErrInsecureBroker stops a hub from handing its password to anyone on the path to a public broker.
var ErrInsecureBroker = errors.New("mqtt: refusing to send credentials in the clear")

// encryptedSchemes are the ones paho wraps in TLS; ws and tcp arrive as plaintext on the wire.
var encryptedSchemes = map[string]bool{
	"ssl": true, "tls": true, "mqtts": true, "mqtt+ssl": true, "tcps": true, "wss": true,
}

// guardTransport allows plaintext only where the traffic cannot leave the machine or the private network.
func guardTransport(brokerURL string, allowInsecure bool) error {
	parsed, err := url.Parse(brokerURL)
	if err != nil {
		return fmt.Errorf("mqtt: parse broker url: %w", err)
	}

	if encryptedSchemes[strings.ToLower(parsed.Scheme)] || allowInsecure {
		return nil
	}

	host := parsed.Hostname()
	if isPrivateHost(host) {
		return nil
	}

	return fmt.Errorf(
		"%w: %s reaches %s over the public internet as plaintext. Use wss:// behind a TLS-terminating proxy, "+
			"or ssl:// against a broker holding certificates. Set MQTT_ALLOW_INSECURE=true to override",
		ErrInsecureBroker, parsed.Scheme+"://", host,
	)
}

// isPrivateHost reports whether traffic to host stays on the loopback, a private range, or an internal name.
func isPrivateHost(host string) bool {
	if host == "" {
		return false
	}

	if ip := net.ParseIP(host); ip != nil {
		return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsUnspecified()
	}

	lower := strings.ToLower(host)
	if lower == "localhost" {
		return true
	}

	// A bare name resolves only inside a container network or a search domain, never on the public internet.
	if !strings.Contains(lower, ".") {
		return true
	}

	for _, suffix := range []string{".internal", ".local", ".localhost", ".localdomain"} {
		if strings.HasSuffix(lower, suffix) {
			return true
		}
	}

	return false
}

// tlsConfigFor returns nil for the system trust store, which already verifies the chain and the hostname.
func tlsConfigFor(caFile string) (*tls.Config, error) {
	if caFile == "" {
		return nil, nil
	}

	pem, err := os.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("mqtt: read ca file: %w", err)
	}

	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM(pem) {
		return nil, fmt.Errorf("mqtt: %s holds no certificates", caFile)
	}

	return &tls.Config{RootCAs: roots, MinVersion: tls.VersionTLS12}, nil
}
