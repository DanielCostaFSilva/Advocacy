package logger

import (
	"context"
	"io"
	"log/slog"
	"os"

	"legalflow/internal/middleware"
)

type Logger struct {
	inner *slog.Logger
}

func New(w io.Writer) *Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	return &Logger{
		inner: slog.New(slog.NewJSONHandler(w, opts)),
	}
}

func NewDefault() *Logger {
	return New(os.Stdout)
}

func (l *Logger) Info(msg string) {
	l.inner.Info(msg)
}

func (l *Logger) Inner() *slog.Logger {
	return l.inner
}

func (l *Logger) InfoContext(ctx context.Context, msg string) {
	attrs := extractAttrs(ctx)
	l.inner.LogAttrs(ctx, slog.LevelInfo, msg, attrs...)
}

func extractAttrs(ctx context.Context) []slog.Attr {
	var attrs []slog.Attr
	if id, ok := ctx.Value(middleware.RequestIDKey).(string); ok {
		attrs = append(attrs, slog.String("request_id", id))
	}
	return attrs
}
