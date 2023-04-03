package main

import (
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/IfedayoAwe/greenlight/internal/jsonlog"
	"github.com/IfedayoAwe/greenlight/internal/mailer"
)

func newTestApplication(t *testing.T) *application {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")
	flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "8502b6a9bc7a9b", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "6d0db8180009fd", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Greenlight <olalekanawe99@gmail.com>", "SMTP sender")
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})
	return &application{
		config: cfg,
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username,
			cfg.smtp.password, cfg.smtp.sender),
		logger: jsonlog.New(io.Discard, jsonlog.LevelInfo),
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)
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

func unMarshal(t *testing.T, js []byte, decodeVar interface{}) {
	js = js[:len(js)-1]

	err := json.Unmarshal(js, &decodeVar)
	if err != nil {
		t.Fatal(err)
	}
}
