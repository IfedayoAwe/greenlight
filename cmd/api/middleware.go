package main

import (
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/IfedayoAwe/greenlight/internal/data"
	"github.com/IfedayoAwe/greenlight/internal/validator"
	"github.com/felixge/httpsnoop"
	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.limiter.enabled {
			ip := realip.FromRequest(r)
			mu.Lock()
			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst),
				}
			}
			clients[ip].lastSeen = time.Now()
			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}
			mu.Unlock()
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]

		v := validator.New()

		if data.ValidateTokenPlaintext(v, token); !v.Valid() {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		r = app.contextSetUser(r, user)

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := app.contextGetUser(r)
		if !user.Activated {
			app.inactiveAccountResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
	return app.requireAuthenticatedUser(fn)
}

func (app *application) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		permissions, err := app.models.Permissions.GetAllForUser(user.ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		if !permissions.Include(code) {
			app.notPermittedResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	}
	return app.requireActivatedUser(fn)
}

func (app *application) requireAdmin(next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if user := app.contextGetUser(r); !user.Admin {
			app.notPermittedResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	}
	return app.requireActivatedUser(fn)
}

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")
		origin := r.Header.Get("Origin")
		if origin != "" && len(app.config.cors.trustedOrigins) != 0 {
			for i := range app.config.cors.trustedOrigins {
				if origin == app.config.cors.trustedOrigins[i] {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					// perflight request
					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
						w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
						w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
						w.WriteHeader(http.StatusOK)
						return
					}

				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

// type WrappedResponseWriter struct {
// 	w  http.ResponseWriter
// 	gw *gzip.Writer
// }

// func NewWrappedResponseWriter(w http.ResponseWriter) *WrappedResponseWriter {
// 	gw := gzip.NewWriter(w)
// 	return &WrappedResponseWriter{w: w, gw: gw}

// }

// func (ww WrappedResponseWriter) Header() http.Header {
// 	return ww.w.Header()
// }

// func (ww WrappedResponseWriter) Write(data []byte) (int, error) {
// 	return ww.gw.Write(data)
// }

// func (ww WrappedResponseWriter) WriteHeader(statuscode int) {
// 	ww.w.WriteHeader(statuscode)
// }

// func (ww WrappedResponseWriter) Flush() {
// 	ww.gw.Flush()
// 	ww.gw.Close()
// }

// func (app *application) enableGzip(next http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Add("Vary", "Accept-Encoding")
// 		acceptEncodingHeader := r.Header.Get("Accept-Encoding")
// 		if strings.Contains(acceptEncodingHeader, "gzip") {
// 			ww := NewWrappedResponseWriter(w)
// 			defer ww.Flush()
// 			ww.Header().Set("Content-Encoding", "gzip")
// 			next.ServeHTTP(ww, r)
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }

func (app *application) metrics(next http.Handler) http.Handler {
	if app.config.metrics.enabled {
		totalRequestsReceived := expvar.NewInt("total_requests_received")
		totalResponsesSent := expvar.NewInt("total_responses_sent")
		totalProcessingTimeMicroseconds := expvar.NewInt("total_processing_time_μs")
		totalProcessingTimeMicrosecondsByMetrics := expvar.NewInt("total_processing_metrics_time_μs")
		totalResponsesSentByStatus := expvar.NewMap("total_responses_sent_by_status")

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			totalRequestsReceived.Add(1)
			// Call the httpsnoop.CaptureMetrics() function, passing in the next handler in
			// the chain along with the existing http.ResponseWriter and http.Request. This
			// returns the metrics struct.
			metrics := httpsnoop.CaptureMetrics(next, w, r)

			totalResponsesSent.Add(1)
			// Calculate the number of microseconds since we began to process the request,
			// then increment the total processing time by this amount.
			duration := time.Since(start).Microseconds()
			totalProcessingTimeMicroseconds.Add(duration)

			// By metrics
			totalProcessingTimeMicrosecondsByMetrics.Add(metrics.Duration.Microseconds())

			// Use the Add() method to increment the count for the given status code by 1.
			// Note that the expvar map is string-keyed, so we need to use the strconv.Itoa()
			// function to convert the status code (which is an integer) to a string.
			totalResponsesSentByStatus.Add(strconv.Itoa(metrics.Code), 1)
		})
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})

}
