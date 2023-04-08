package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"

	"github.com/IfedayoAwe/greenlight/internal/data/mock"
	"github.com/IfedayoAwe/greenlight/internal/jsonlog"
	"github.com/IfedayoAwe/greenlight/internal/mailer"
)

func newTestApplication(t *testing.T) *application {
	var testCfg config
	testCfg.env = "development"
	testCfg.limiter.rps = 2
	testCfg.limiter.burst = 4
	testCfg.limiter.enabled = false
	testCfg.smtp.host = "sandbox.smtp.mailtrap.io"
	testCfg.smtp.port = 2525
	testCfg.smtp.username = "8502b6a9bc7a9b"
	testCfg.smtp.password = "6d0db8180009fd"
	testCfg.smtp.sender = "Greenlight <olalekanawe99@gmail.com>"
	testCfg.cors.trustedOrigins = []string{"*"}

	return &application{
		config: testCfg,
		logger: jsonlog.New(io.Discard, jsonlog.LevelInfo),
		models: mock.NewMockModels(),
		mailer: mailer.New(testCfg.smtp.host, testCfg.smtp.port, testCfg.smtp.username, testCfg.smtp.password, testCfg.smtp.sender),
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	ts.Client().Jar = jar
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &testServer{ts}
}

func (ts *testServer) do(t *testing.T, r *http.Request) (int, http.Header, []byte) {
	rs, err := ts.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	return rs.StatusCode, rs.Header, body
}

// func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, []byte) {
// 	rs, err := ts.Client().PostForm(ts.URL+urlPath, form)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer rs.Body.Close()
// 	body, err := io.ReadAll(rs.Body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	return rs.StatusCode, rs.Header, body
// }

func unMarshal(t *testing.T, js []byte, decodeVar interface{}) {
	js = js[:len(js)-1]

	err := json.Unmarshal(js, &decodeVar)
	if err != nil {
		t.Fatal(err)
	}
}
