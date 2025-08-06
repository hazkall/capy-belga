package middlewares

import (
	"context"
	"net/http"

	"github.com/hazkall/capy-belga/pkg/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func RequestsCountMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		telemetry.RequestCounter.Add(context.Background(), 1,
			metric.WithAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.path", r.URL.Path),
			))
		next.ServeHTTP(w, r)
	})
}
