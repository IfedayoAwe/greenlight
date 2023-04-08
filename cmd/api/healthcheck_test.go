package main

import (
	"net/http"
	"reflect"
	"testing"
)

func TestHealthcheckHandler(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	env := envelope{
		"status": "available",
		"system_info": map[string]interface{}{
			"environment": app.config.env,
			"version":     version,
		},
	}

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/v1/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	code, header, body := ts.do(t, req)

	var decVar envelope

	unMarshal(t, body, &decVar)

	contentType := header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("want %q; got %q", "application/json", contentType)
	}

	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	if !reflect.DeepEqual(env, decVar) {
		t.Errorf("want body to equal %q", env)
	}

}
