package event_test

import (
	"testing"
	"time"

	"ava/api/internal/services/event"

	"github.com/google/uuid"
)

func TestSubscribersOnlySeeTheirOwnTenant(t *testing.T) {
	broker := event.NewService()

	mine, cancelMine := broker.Subscribe(uuid.MustParse("11111111-1111-1111-1111-111111111111"))
	defer cancelMine()

	theirs, cancelTheirs := broker.Subscribe(uuid.MustParse("22222222-2222-2222-2222-222222222222"))
	defer cancelTheirs()

	broker.Publish(uuid.MustParse("11111111-1111-1111-1111-111111111111"), []byte("for me"))

	select {
	case got := <-mine:
		if string(got) != "for me" {
			t.Errorf("got %q", got)
		}
	case <-time.After(time.Second):
		t.Fatal("did not receive our own tenant's event")
	}

	select {
	case leaked := <-theirs:
		t.Errorf("another tenant received %q", leaked)
	case <-time.After(100 * time.Millisecond):
	}
}

func TestPublishDoesNotBlockOnASlowClient(t *testing.T) {
	broker := event.NewService()
	tenant := uuid.New()

	_, cancel := broker.Subscribe(tenant)
	defer cancel()

	done := make(chan struct{})

	go func() {
		for range 1000 {
			broker.Publish(tenant, []byte("tick"))
		}

		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("a client that never reads blocked the publisher")
	}
}

func TestCancelRemovesTheSubscriber(t *testing.T) {
	broker := event.NewService()
	tenant := uuid.New()

	_, cancel := broker.Subscribe(tenant)

	if broker.Clients(tenant) != 1 {
		t.Fatalf("clients = %d", broker.Clients(tenant))
	}

	cancel()
	cancel()

	if broker.Clients(tenant) != 0 {
		t.Errorf("clients after cancel = %d", broker.Clients(tenant))
	}
}
