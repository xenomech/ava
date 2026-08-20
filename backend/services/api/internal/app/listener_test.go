package app

import "testing"

func TestHubIsParsedFromEitherTopic(t *testing.T) {
	cases := []struct {
		name  string
		topic string
		want  string
	}{
		{
			name:  "state topic",
			topic: "ava/acme/6f1e2c4a-0000-4000-8000-000000000001/state",
			want:  "6f1e2c4a-0000-4000-8000-000000000001",
		},
		{
			name:  "status topic",
			topic: "ava/acme/6f1e2c4a-0000-4000-8000-000000000001/status",
			want:  "6f1e2c4a-0000-4000-8000-000000000001",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			hubID, ok := hubFromTopic(tc.topic)
			if !ok {
				t.Fatalf("did not parse %q", tc.topic)
			}

			if hubID.String() != tc.want {
				t.Errorf("got %s, want %s", hubID, tc.want)
			}
		})
	}
}

func TestATopicWithoutAHubIsIgnored(t *testing.T) {
	for _, topic := range []string{"ava/acme/state", "ava/acme/not-a-uuid/status", "ava/acme/x/y/status"} {
		if _, ok := hubFromTopic(topic); ok {
			t.Errorf("accepted %q", topic)
		}
	}
}
