package linkding

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientDoSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Token secret-token" {
			t.Fatalf("bad auth: %q", r.Header.Get("Authorization"))
		}
		if r.URL.Path != "/api/bookmarks/" {
			t.Fatalf("path %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"count":1,"results":[]}`))
	}))
	defer srv.Close()

	c := &Client{Token: "secret-token", BaseURL: srv.URL, HTTP: srv.Client()}
	raw, err := c.Do(context.Background(), Request{Method: http.MethodGet, Path: "/api/bookmarks/"})
	if err != nil {
		t.Fatal(err)
	}
	var page map[string]any
	if err := json.Unmarshal(raw, &page); err != nil || page["count"].(float64) != 1 {
		t.Fatalf("got %s", raw)
	}
}

func TestClientRedact(t *testing.T) {
	c := &Client{Token: "abc123"}
	if got := c.Redact("Token abc123 leaked"); got != "Token <TOKEN> leaked" {
		t.Fatalf("redact: %q", got)
	}
}

func TestDryRunMasksToken(t *testing.T) {
	d := DryRun("https://ld.test", Request{Method: "GET", Path: "/api/tags/"})
	h := d["header"].(map[string]string)
	if h["Authorization"] != "Token <TOKEN>" {
		t.Fatalf("header %v", h)
	}
}

func TestClientDoAuthError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"detail":"Invalid token."}`))
	}))
	defer srv.Close()

	c := &Client{Token: "bad", BaseURL: srv.URL, HTTP: srv.Client()}
	_, err := c.Do(context.Background(), Request{Method: http.MethodGet, Path: "/api/user/profile/"})
	var api *APIError
	if !errors.As(err, &api) || !api.IsAuth() {
		t.Fatalf("expected auth error, got %v", err)
	}
}
