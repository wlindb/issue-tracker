package api

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type logLineKey struct{}

// logLine accumulates structured fields over the lifetime of a single request.
// It is safe for concurrent use by handlers that spawn goroutines.
type logLine struct {
	mu     sync.Mutex
	fields []slog.Attr
}

func (l *logLine) add(attrs ...slog.Attr) {
	l.mu.Lock()
	l.fields = append(l.fields, attrs...)
	l.mu.Unlock()
}

// AddLogFields attaches structured fields to the canonical log line for the
// current request. Call this from handlers or services that receive the context.
func AddLogFields(ctx context.Context, attrs ...slog.Attr) {
	if ll, ok := ctx.Value(logLineKey{}).(*logLine); ok {
		ll.add(attrs...)
	}
}

// RequestLogger returns an Echo middleware that emits one structured log line
// per request containing HTTP metadata and any fields added via AddLogFields.
func RequestLogger(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			ll := &logLine{}

			req := c.Request()
			ctx := context.WithValue(req.Context(), logLineKey{}, ll)
			c.SetRequest(req.WithContext(ctx))

			err := next(c)

			status := c.Response().Status
			if err != nil {
				if he, ok := err.(*echo.HTTPError); ok {
					status = he.Code
				} else {
					status = http.StatusInternalServerError
				}
			}

			attrs := []slog.Attr{
				slog.String("method", req.Method),
				slog.String("path", req.URL.Path),
				slog.Int("status", status),
				slog.Int64("duration_ms", time.Since(start).Milliseconds()),
				slog.Int64("bytes_out", c.Response().Size),
			}
			if rid := c.Response().Header().Get(echo.HeaderXRequestID); rid != "" {
				attrs = append(attrs, slog.String("request_id", rid))
			}
			if err != nil {
				attrs = append(attrs, slog.String("error", err.Error()))
			}

			ll.mu.Lock()
			attrs = append(attrs, ll.fields...)
			ll.mu.Unlock()

			level := slog.LevelInfo
			if status >= 500 {
				level = slog.LevelError
			} else if status >= 400 {
				level = slog.LevelWarn
			}

			logger.LogAttrs(ctx, level, "request", attrs...)

			return err
		}
	}
}
