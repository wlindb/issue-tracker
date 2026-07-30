package gemini_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/genai"

	"github.com/wlindb/issue-tracker/internal/infrastructure/search/gemini"
)

func embeddingServer(t *testing.T, statusCode int, values []float32) *httptest.Server {
	t.Helper()
	body, err := json.Marshal(map[string]any{
		"embeddings": []map[string]any{
			{"values": values},
		},
	})
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write(body)
	}))
	t.Cleanup(server.Close)
	return server
}

func newTestGenerator(t *testing.T, server *httptest.Server) *gemini.GeminiEmbeddingGenerator {
	t.Helper()
	generator, err := gemini.NewGeminiEmbeddingGenerator(
		context.Background(),
		"fake-api-key",
		"test-model",
		gemini.WithHTTPOptions(genai.HTTPOptions{BaseURL: server.URL + "/"}),
	)
	require.NoError(t, err)
	return generator
}

func Test_GenerateEmbedding_TitleAndDescription_ReturnsVector(t *testing.T) {
	expected := []float32{0.1, 0.2, 0.3, 0.4}
	server := embeddingServer(t, http.StatusOK, expected)

	generator := newTestGenerator(t, server)

	actual, err := generator.GenerateEmbedding(context.Background(), "Test Title", "Test description body")

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func Test_GenerateEmbedding_EmptyTitleAndDescription_ReturnsError(t *testing.T) {
	server := embeddingServer(t, http.StatusOK, nil)

	generator := newTestGenerator(t, server)

	_, err := generator.GenerateEmbedding(context.Background(), "", "")

	require.Error(t, err)
	assert.True(t, errors.Is(err, gemini.ErrEmptyContent))
}

func Test_GenerateEmbedding_ServerError_ReturnsError(t *testing.T) {
	server := embeddingServer(t, http.StatusInternalServerError, nil)

	generator := newTestGenerator(t, server)

	_, err := generator.GenerateEmbedding(context.Background(), "Test Title", "Test description body")

	require.Error(t, err)
}

func Test_GenerateQueryEmbedding_ValidQuery_ReturnsVector(t *testing.T) {
	expected := []float32{0.1, 0.2, 0.3, 0.4}
	server := embeddingServer(t, http.StatusOK, expected)

	generator := newTestGenerator(t, server)

	actual, err := generator.GenerateQueryEmbedding(context.Background(), "Test query")

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func Test_GenerateQueryEmbedding_EmptyQuery_ReturnsError(t *testing.T) {
	server := embeddingServer(t, http.StatusOK, nil)

	generator := newTestGenerator(t, server)

	_, err := generator.GenerateQueryEmbedding(context.Background(), "")

	require.Error(t, err)
	assert.True(t, errors.Is(err, gemini.ErrEmptyQuery))
}

func Test_GenerateQueryEmbedding_ServerError_ReturnsError(t *testing.T) {
	server := embeddingServer(t, http.StatusInternalServerError, nil)

	generator := newTestGenerator(t, server)

	_, err := generator.GenerateQueryEmbedding(context.Background(), "Test query")

	require.Error(t, err)
}
