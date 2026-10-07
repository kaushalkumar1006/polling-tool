package app

import (
	"net/http"
	"testing"
	"time"
)

func TestProductionSessionCookieIsHostOnly(t *testing.T) {
	a := &App{cfg: Config{CookieSecure: true}}
	cookie := a.sessionCookie("token", 86400)

	if cookie.Name != "__Host-session" || cookie.Domain != "" || cookie.Path != "/" {
		t.Fatalf("expected a host-only production cookie, got %#v", cookie)
	}
	if !cookie.Secure || !cookie.HttpOnly || cookie.SameSite != http.SameSiteNoneMode {
		t.Fatalf("expected Secure, HttpOnly, SameSite=None, got %#v", cookie)
	}
}

func TestValidatePollRejectsDuplicateOptions(t *testing.T) {
	_, err := validatePoll("Choose a direction", []string{"React", "react"}, nil)
	if err == nil { t.Fatal("expected duplicate options to be rejected") }
}

func TestValidatePollAcceptsWellFormedPoll(t *testing.T) {
	options, err := validatePoll("Where should we meet?", []string{"Studio", "Courtyard"}, nil)
	if err != nil { t.Fatalf("unexpected validation error: %v", err) }
	if len(options) != 2 || options[0].ID == "" { t.Fatalf("expected generated option IDs, got %#v", options) }
}

func TestValidatePollRejectsInvalidShape(t *testing.T) {
	for _, input := range [][]string{{"Only one"}, {"A", "A", "A"}} {
		if _, err := validatePoll("Valid question", input, nil); err == nil { t.Fatalf("expected invalid options %v to fail", input) }
	}
}

func TestValidateCredentialsRejectsInvalidInput(t *testing.T) {
	if validateCredentials("not-an-email", "password123") == nil {
		t.Fatal("expected invalid email to fail")
	}
	if validateCredentials("person@example.com", "short") == nil {
		t.Fatal("expected short password to fail")
	}
}

func TestValidatePollRejectsPastExpiration(t *testing.T) {
	past := time.Now().Add(-time.Minute)
	if _, err := validatePoll("A valid question", []string{"One", "Two"}, &past); err == nil {
		t.Fatal("expected past expiration to fail")
	}
}
