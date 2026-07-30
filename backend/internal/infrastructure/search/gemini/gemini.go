package gemini

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/genai"
)

// ErrEmptyContent is returned when both title and description are empty.
var ErrEmptyContent = errors.New("title and description are both empty")

// EmbeddingGenerator generates vector embeddings for issue content.
type EmbeddingGenerator interface {
	GenerateEmbedding(ctx context.Context, title, description string) ([]float32, error)
}

// GeminiEmbeddingGenerator generates embeddings using the Google Gemini API.
type GeminiEmbeddingGenerator struct {
	client *genai.Client
	model  string
}

// Option is a functional option for configuring a GeminiEmbeddingGenerator.
type Option func(*genai.ClientConfig)

// WithHTTPOptions returns an Option that sets the HTTP options on the Gemini client config.
// This is primarily useful in tests to redirect requests to a local server.
func WithHTTPOptions(httpOptions genai.HTTPOptions) Option {
	return func(config *genai.ClientConfig) {
		config.HTTPOptions = httpOptions
	}
}

// NewGeminiEmbeddingGenerator creates a GeminiEmbeddingGenerator that calls the Gemini API
// with the given API key and embedding model.
func NewGeminiEmbeddingGenerator(ctx context.Context, apiKey, model string, opts ...Option) (*GeminiEmbeddingGenerator, error) {
	config := &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	}
	for _, opt := range opts {
		opt(config)
	}

	client, err := genai.NewClient(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create gemini client: %w", err)
	}

	return &GeminiEmbeddingGenerator{
		client: client,
		model:  model,
	}, nil
}

// GenerateEmbedding returns a vector embedding for the given issue title and description.
// It returns ErrEmptyContent if both title and description are empty.
func (g *GeminiEmbeddingGenerator) GenerateEmbedding(ctx context.Context, title, description string) ([]float32, error) {
	var zero []float32

	if title == "" && description == "" {
		return zero, ErrEmptyContent
	}

	combined := buildContent(title, description)

	result, err := g.client.Models.EmbedContent(ctx, g.model, genai.Text(combined), &genai.EmbedContentConfig{
		TaskType: "RETRIEVAL_DOCUMENT",
	})
	if err != nil {
		return zero, fmt.Errorf("generate embedding: %w", err)
	}

	if len(result.Embeddings) == 0 {
		return zero, fmt.Errorf("generate embedding: no embeddings returned")
	}

	return result.Embeddings[0].Values, nil
}

// buildContent combines title and description into a single plain-text string for embedding,
// without field labels, so the content reads as natural prose.
func buildContent(title, description string) string {
	switch {
	case title == "":
		return description
	case description == "":
		return title
	default:
		return title + "\n\n" + description
	}
}
