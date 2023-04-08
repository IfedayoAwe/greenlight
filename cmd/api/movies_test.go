package main

import (
	"bytes"
	"net/http"
	"testing"
)

func TestShowMovies(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
		token    string
	}{
		{"Unauthenticated", "/v1/movies/1", http.StatusUnauthorized, nil, ""},
		{"Authenticated", "/v1/movies/1", http.StatusOK, []byte("Test Movie"), "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRI"},
		{"UnActivated", "/v1/movies/1", http.StatusForbidden, nil, "Bearer HTE34GKUHNDUSJ3QRUT6IKWKRJ"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req, err := http.NewRequest(http.MethodGet, ts.URL+tt.urlPath, nil)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Authorization", tt.token)

			code, header, body := ts.do(t, req)
			contentType := header.Get("Content-Type")
			if contentType != "application/json" {
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
