package broker

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestProvisionAgainstALiveBroker(t *testing.T) {
	url := os.Getenv("AVA_TEST_BROKER")
	if url == "" {
		t.Skip("set AVA_TEST_BROKER, AVA_TEST_BROKER_USER and AVA_TEST_BROKER_PASS to run this against a real broker")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	b, err := Connect(ctx, Config{
		URL: url, ClientID: "provisioner",
		Username: os.Getenv("AVA_TEST_BROKER_USER"), Password: os.Getenv("AVA_TEST_BROKER_PASS"),
	})
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer b.Close()

	for _, pair := range [][2]string{{"acme", "hub-a"}, {"globex", "hub-b"}} {
		username, password, err := b.ProvisionHub(ctx, pair[0], pair[1])
		if err != nil {
			t.Fatalf("provision %s: %v", pair[1], err)
		}

		t.Logf("CREDS %s %s", username, password)
	}

	time.Sleep(2 * time.Second)
}
