package middlewares

import (
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/hazkall/capy-belga/pkg/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *statusResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func OtelMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			remoteIP = r.RemoteAddr
		}

		startTime := time.Now()

		ctx, span := telemetry.Tracer.Start(r.Context(), r.URL.Path,
			trace.WithAttributes(
				attribute.String("http.server.remote_ip", remoteIP),
				attribute.String("http.server.protocol", r.Proto),
				attribute.String("http.server.host", r.Host),
				attribute.String("http.server.path", r.URL.Path),
				attribute.String("http.server.method", r.Method),
				attribute.String("http.server.user_agent", r.UserAgent()),
				attribute.String("http.server.referer", r.Referer()),
				attribute.String("http.server.request_uri", r.RequestURI),
				attribute.String("http.server.request_scheme", r.URL.Scheme),
			),
		)

		defer span.End()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w}
		next.ServeHTTP(wrappedWriter, r.WithContext(ctx))

		if wrappedWriter.statusCode >= 200 && wrappedWriter.statusCode < 400 {
			span.SetStatus(codes.Ok, "HTTP Success")
		} else {
			span.SetStatus(codes.Error, "HTTP Error")
		}

		duration := time.Since(startTime)
		statusCode := wrappedWriter.statusCode

		span.SetAttributes(
			attribute.Int("http.server.status_code", statusCode),
			attribute.Float64("http.server.duration_ms", float64(duration.Milliseconds())),
		)

		slog.InfoContext(ctx, "Request completed",
			slog.String("trace.id", span.SpanContext().TraceID().String()),
			slog.String("trace.span_id", span.SpanContext().SpanID().String()),
			slog.String("http.server.remote_ip", remoteIP),
			slog.String("http.server.protocol", r.Proto),
			slog.String("http.server.host", r.Host),
			slog.String("http.server.path", r.URL.Path),
			slog.String("http.server.method", r.Method),
			slog.String("http.server.user_agent", r.UserAgent()),
			slog.String("http.server.referer", r.Referer()),
			slog.String("http.server.request_uri", r.RequestURI),
			slog.String("http.server.request_scheme", r.URL.Scheme),
			slog.Int("http.server.status_code", statusCode),
			slog.Float64("http.server.duration_ms", float64(duration.Milliseconds())),
		)
	})

}
