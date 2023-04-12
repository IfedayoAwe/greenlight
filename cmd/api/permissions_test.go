package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestAddMovieWritePermission(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		email    string
		wantCode int
		wantBody []byte
		token    string
	}{
		{"Created", "ayo@gmail.com", http.StatusOK, []byte("ayo@gmail.com"), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI"},
		{"InvalidEmail", "gmail.com", http.StatusUnprocessableEntity, []byte("must be a valid email address"), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI"},
		{"UnknownEmail", "foo@gmail.com", http.StatusUnprocessableEntity, []byte("no matching email address found"), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI"},
		{"InactiveAccount", "mummy@gmail.com", http.StatusUnprocessableEntity, []byte("This user already has this permission"), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI"},
		{"InvalidToken", "mummy@gmail.com", http.StatusUnauthorized, []byte("invalid or missing authentication token"), "HTE34GKUHNDUSJ3QRUT6IKWKRI"},
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

			req, err := http.NewRequest(http.MethodPost, ts.URL+"/v1/users/movie-permission", bytes.NewReader(payload))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Authorization", tt.token)
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
