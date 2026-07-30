-- +goose Up
CREATE EXTENSION IF NOT EXISTS vector;
ALTER TABLE issue_documents ADD COLUMN embedding vector(1536);
CREATE INDEX issue_documents_embedding_hnsw_idx ON issue_documents USING hnsw (embedding vector_cosine_ops);

-- +goose Down
DROP INDEX IF EXISTS issue_documents_embedding_hnsw_idx;
ALTER TABLE issue_documents DROP COLUMN embedding;
DROP EXTENSION IF EXISTS vector;
