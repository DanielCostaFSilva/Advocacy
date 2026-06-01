package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"legalflow/internal/middleware"
)

func TestLogger_ShouldOutputJSON(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf)

	log.Info("hello world")

	var entry map[string]any
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Equal(t, "hello world", entry["msg"])
	assert.Equal(t, "INFO", entry["level"])
	assert.NotEmpty(t, entry["time"])
}

func TestLogger_ShouldIncludeRequestID(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf)

	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "req-123")
	log.InfoContext(ctx, "with request")

	var entry map[string]any
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Equal(t, "req-123", entry["request_id"])
}
