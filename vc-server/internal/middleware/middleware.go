package server

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
	"varchess/internal/logger"

	"github.com/go-chi/cors"
	"github.com/google/uuid"
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


func RequestLogger(l logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestID, correlationID := getOrCreateIDs(r)
			reqLogger:= l.With().Str("requestID", requestID).Str("correlationID", correlationID).Logger()
			ctx := r.Context()
			ctx = context.WithValue(ctx,logger.LoggerKey,reqLogger)
			r = r.WithContext(ctx)
			lrw := newLoggingResponseWriter(w)
			
			defer func() {
				logLevel := zerolog.InfoLevel
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
				queryString := ""
				if r.URL.RawQuery != "" {
					queryString = "?" + r.URL.RawQuery
				}

				reqLogger.WithLevel(logLevel).
					Dur("duration", time.Since(start)).
					Str("info", fmt.Sprintf("[%v] %s: %s%s", lrw.statusCode, r.Method, r.URL.RequestURI(), queryString)).Msg("Request processed")
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

func getOrCreateIDs(r *http.Request) (reqID string, corrID string) {
	reqID = getRequestID(r)
	corrID = getCorrelationID(r)
	if reqID == "" {
		reqID = uuid.NewString()
	}
	if corrID == "" {
		corrID = uuid.NewString()
	}
	return reqID, corrID
}

func getRequestID(r *http.Request) string {
	return r.Header.Get("X-Request-ID")
}

func getCorrelationID(r *http.Request) string {
	return r.Header.Get("X-Correlation-id")
}