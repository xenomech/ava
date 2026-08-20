package event

import (
	"encoding/json"
	"sync"

	"ava/pkg/logger"

	"github.com/google/uuid"
)

const bufferPerClient = 16

type Service interface {
	Subscribe(tenantID uuid.UUID) (events <-chan []byte, cancel func())
	Publish(tenantID uuid.UUID, payload []byte)
	PublishJSON(tenantID uuid.UUID, event any)
	Clients(tenantID uuid.UUID) int
}

type broker struct {
	mu      sync.RWMutex
	tenants map[uuid.UUID]map[chan []byte]struct{}
}

func NewService() Service {
	return &broker{tenants: make(map[uuid.UUID]map[chan []byte]struct{})}
}

func (b *broker) Subscribe(tenantID uuid.UUID) (events <-chan []byte, cancel func()) {
	stream := make(chan []byte, bufferPerClient)

	b.mu.Lock()

	if b.tenants[tenantID] == nil {
		b.tenants[tenantID] = make(map[chan []byte]struct{})
	}

	b.tenants[tenantID][stream] = struct{}{}
	b.mu.Unlock()

	var once sync.Once

	cancel = func() {
		once.Do(func() {
			b.mu.Lock()
			defer b.mu.Unlock()

			if listeners, ok := b.tenants[tenantID]; ok {
				delete(listeners, stream)

				if len(listeners) == 0 {
					delete(b.tenants, tenantID)
				}
			}

			close(stream)
		})
	}

	return stream, cancel
}

func (b *broker) Publish(tenantID uuid.UUID, payload []byte) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for stream := range b.tenants[tenantID] {
		select {
		case stream <- payload:
		default:
		}
	}
}

func (b *broker) PublishJSON(tenantID uuid.UUID, event any) {
	payload, err := json.Marshal(event)
	if err != nil {
		logger.Warn("EVENT_ENCODE_FAILED", logger.Err(err))

		return
	}

	b.Publish(tenantID, payload)
}

func (b *broker) Clients(tenantID uuid.UUID) int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return len(b.tenants[tenantID])
}
