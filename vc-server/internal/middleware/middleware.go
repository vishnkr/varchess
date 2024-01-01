package server

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
)


type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}


func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) Hijack()(net.Conn, *bufio.ReadWriter, error) {
    h, ok := lrw.ResponseWriter.(http.Hijacker)
    if !ok {
        return nil, nil, errors.New("hijack not supported")
    }
    return h.Hijack()
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}


func RequestLogger(l zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler{
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
}

func Cors() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	})
}

