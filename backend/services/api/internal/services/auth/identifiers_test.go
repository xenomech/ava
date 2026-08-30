package auth

import "testing"

func TestIdentifierBase(t *testing.T) {
	cases := map[string]string{
		"gokul@runable.com":            "gokul",
		"Gokul.Nath+lights@Home.co.uk": "gokul-nath-lights",
		"a@b.com":                      identifierFallback,
		"____@b.com":                   identifierFallback,
		"日本@example.com":               identifierFallback,
		"averyveryverylongmailboxname@example.com": "averyveryverylongmailbox",
	}

	for email, want := range cases {
		if got := identifierBase(email); got != want {
			t.Errorf("identifierBase(%q) = %q, want %q", email, got, want)
		}
	}
}

// The base must survive truncation as a valid slug, with no trailing separator from a cut on a hyphen.
func TestIdentifierBaseTruncatesToAValidSlug(t *testing.T) {
	got := identifierBase("aaaaaaaaaaaaaaaaaaaaaaa-bbbbbbb@example.com")

	if got != "aaaaaaaaaaaaaaaaaaaaaaa" {
		t.Fatalf("identifierBase = %q", got)
	}
}

func TestFirstFreeSkipsTakenNames(t *testing.T) {
	taken := map[string]bool{"gokul": true, "gokul-2": true}

	got, err := firstFree("gokul", func(name string) (bool, error) {
		return !taken[name], nil
	})
	if err != nil {
		t.Fatalf("firstFree: %v", err)
	}

	if got != "gokul-3" {
		t.Fatalf("firstFree = %q, want gokul-3", got)
	}
}

// Every candidate being taken must still yield a name rather than refusing the registration.
func TestFirstFreeFallsBackToARandomSuffix(t *testing.T) {
	got, err := firstFree("gokul", func(string) (bool, error) { return false, nil })
	if err != nil {
		t.Fatalf("firstFree: %v", err)
	}

	if len(got) != len("gokul-")+8 {
		t.Fatalf("firstFree = %q, want a random suffix", got)
	}
}
