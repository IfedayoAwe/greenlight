package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestCreateAuthenticationToken(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		email    string
		password string
		wantCode int
		wantBody []byte
	}{
		{"Created", "olalekanawe99@gmail.com", "1234567890", http.StatusCreated, []byte("token")},
		{"InvalidEmail", "ola.com", "1234567890", http.StatusUnprocessableEntity, []byte("must be a valid email address")},
		{"EmptyEmail", "", "1234567890", http.StatusUnprocessableEntity, []byte("\"email\": \"must be provided\"")},
		{"UnknownEmail", "ola99@gmail.com", "1234567890", http.StatusUnauthorized, []byte("invalid authentication credentials")},
		{"UnactivatedAccount", "ayo@gmail.com", "1234567890", http.StatusForbidden, []byte("your user account must be activated to access this resource")},
		{"Created", "olalekanawe99@gmail.com", "12345678904", http.StatusUnauthorized, []byte("invalid authentication credentials")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			credentials := struct {
				Email    string
				Password string
			}{
				Email:    tt.email,
				Password: tt.password,
			}

			payload, err := json.Marshal(credentials)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, ts.URL+"/v1/tokens/authentication", bytes.NewReader(payload))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Add("Content-Type", "application/json")

			code, header, body := ts.do(t, req)
			if contentType := header.Get("Content-Type"); contentType != "application/json" {
				t.Errorf("want %q; got %q", "application/json", contentType)
			}

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q", tt.wantBody)
			}
		})
	}
}

func TestCreateActivationToken(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		email    string
		wantCode int
		wantBody []byte
	}{
		{"Created", "ayo@gmail.com", http.StatusAccepted, []byte("an email will be sent to you containing activation instructions")},
		{"AlreadyActivated", "olalekanawe99@gmail.com", http.StatusUnprocessableEntity, []byte("user has already been activated")},
		{"InvalidEmail", "gmail.com", http.StatusUnprocessableEntity, []byte("must be a valid email address")},
		{"UnknownEmail", "foo@gmail.com", http.StatusUnprocessableEntity, []byte("no matching email address found")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			credentials := struct {
				Email string
			}{
				Email: tt.email,
			}

			payload, err := json.Marshal(credentials)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, ts.URL+"/v1/tokens/activation", bytes.NewReader(payload))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Add("Content-Type", "application/json")

			code, header, body := ts.do(t, req)
			if contentType := header.Get("Content-Type"); contentType != "application/json" {
				t.Errorf("want %q; got %q", "application/json", contentType)
			}

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q", tt.wantBody)
			}
		})
	}
}
