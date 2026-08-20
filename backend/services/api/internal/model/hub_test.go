package model_test

import (
	"testing"
	"time"

	"ava/api/internal/model"
)

func TestPresenceNeedsBothTheFlagAndRecentContact(t *testing.T) {
	fresh := time.Now().Add(-time.Minute)
	stale := time.Now().Add(-model.PresenceGrace - time.Minute)

	cases := []struct {
		name       string
		online     bool
		lastSeenAt *time.Time
		want       bool
	}{
		{name: "connected and heard from recently", online: true, lastSeenAt: &fresh, want: true},
		{name: "connected but silent for too long", online: true, lastSeenAt: &stale, want: false},
		{name: "disconnected but heard from recently", online: false, lastSeenAt: &fresh, want: false},
		{name: "never heard from", online: true, lastSeenAt: nil, want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			hub := &model.Hub{Online: tc.online, LastSeenAt: tc.lastSeenAt}

			if got := hub.IsOnline(); got != tc.want {
				t.Errorf("IsOnline() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestPresenceIsUnknownUntilTheHubIsHeardFrom(t *testing.T) {
	now := time.Now()

	if (&model.Hub{}).PresenceKnown() {
		t.Error("a hub that has never reported should have unknown presence")
	}

	if !(&model.Hub{LastSeenAt: &now}).PresenceKnown() {
		t.Error("a hub that has reported should have known presence")
	}
}
