package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/IfedayoAwe/greenlight/internal/data/mock"
	"github.com/IfedayoAwe/greenlight/internal/jsonlog"
	"github.com/IfedayoAwe/greenlight/internal/mailer"
)

func newTestApplication(t *testing.T) *application {
	var testCfg config
	// // flag.IntVar(&testCfg.port, "port", 4000, "API server port")
	// flag.StringVar(&testCfg.env, "env", "development", "Environment (development|staging|production)")
	// flag.StringVar(&testCfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")
	// flag.IntVar(&testCfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	// flag.IntVar(&testCfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	// flag.StringVar(&testCfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	// flag.Float64Var(&testCfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	// flag.IntVar(&testCfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	// flag.BoolVar(&testCfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")
	// flag.StringVar(&testCfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	// flag.IntVar(&testCfg.smtp.port, "smtp-port", 2525, "SMTP port")
	// flag.StringVar(&testCfg.smtp.username, "smtp-username", "8502b6a9bc7a9b", "SMTP username")
	// flag.StringVar(&testCfg.smtp.password, "smtp-password", "6d0db8180009fd", "SMTP password")
	// flag.StringVar(&testCfg.smtp.sender, "smtp-sender", "Greenlight <olalekanawe99@gmail.com>", "SMTP sender")
	// flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
	// 	testCfg.cors.trustedOrigins = strings.Fields(val)
	// 	return nil
	// })
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

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
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

func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, []byte) {
	rs, err := ts.Client().PostForm(ts.URL+urlPath, form)
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

func unMarshal(t *testing.T, js []byte, decodeVar interface{}) {
	js = js[:len(js)-1]

	err := json.Unmarshal(js, &decodeVar)
	if err != nil {
		t.Fatal(err)
	}
}
