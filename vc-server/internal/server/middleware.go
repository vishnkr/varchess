package server

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := newLoggingResponseWriter(w)
		logLevel := zerolog.InfoLevel

		defer func() {
			panicVal := recover()
			if panicVal != nil {
				lrw.statusCode = http.StatusInternalServerError
				panic(panicVal)
			}

			if lrw.statusCode >= http.StatusInternalServerError {
				logLevel = zerolog.ErrorLevel
			} else if lrw.statusCode >= http.StatusBadRequest {
				logLevel = zerolog.WarnLevel
			}
			l.
				WithLevel(logLevel).
				Str("method", r.Method).
				Str("url", r.URL.RequestURI()).
				Dur("elapsed_ms", time.Since(start)).
				Int("status_code", lrw.statusCode).
				Msg("incoming request")
		}()
		next.ServeHTTP(lrw, r)
	})
}

func AllowOptions(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {

			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
			w.Header().Set("Access-Control-Max-Age", "3600")
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
