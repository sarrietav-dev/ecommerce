package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/google/jsonapi"
	"github.com/sarrietav-dev/ecommerce/catalog/internal"
	"github.com/sarrietav-dev/ecommerce/catalog/internal/logger"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger.Logger.Info("Incoming request",
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
			slog.String("remote_addr", r.RemoteAddr),
		)

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log request processing time
		logger.Logger.Info("Request processed",
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
			slog.Duration("duration", time.Since(start)),
		)
	})
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Logger.Error("Unhandled error", slog.Any("error", err))
				w.Header().Set("Content-Type", jsonapi.MediaType)
				w.WriteHeader(http.StatusInternalServerError)
				internal.WriteErrorResponse(w, err.(error), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
