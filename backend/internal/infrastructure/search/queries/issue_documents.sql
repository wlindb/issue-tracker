-- name: CreateIssueDocument :one
INSERT INTO issue_documents (id, workspace_id, title, description, embedding, created_at, updated_at)
VALUES (@id, current_setting('app.workspace_id')::uuid, @title, @description, @embedding, NOW(), NOW())
ON CONFLICT (id) DO NOTHING
RETURNING *;

-- name: GetIssueDocument :one
SELECT * FROM issue_documents WHERE id = @id;

-- name: ListIssueDocuments :many
SELECT * FROM issue_documents
WHERE workspace_id = @workspace_id
ORDER BY created_at DESC;

-- name: UpdateIssueDocument :one
UPDATE issue_documents
SET title       = @title,
    description = @description,
    updated_at  = NOW()
WHERE id = @id
  AND updated_at = @updated_at
RETURNING *;

-- name: FindIssueDocumentsByDescription :many
WITH vector_search AS (
    SELECT id,
           ROW_NUMBER() OVER (ORDER BY embedding <=> @embedding::vector) AS rank
    FROM issue_documents
    WHERE workspace_id = current_setting('app.workspace_id')::uuid
      AND embedding IS NOT NULL
    ORDER BY embedding <=> @embedding::vector
    LIMIT 20
),
keyword_search AS (
    SELECT id,
           ROW_NUMBER() OVER (ORDER BY description <@> to_bm25query(@description::text, 'issue_documents_description_bm25_idx')) AS rank
    FROM issue_documents
    WHERE workspace_id = current_setting('app.workspace_id')::uuid
      AND description <@> to_bm25query(@description::text, 'issue_documents_description_bm25_idx') < 0
    ORDER BY description <@> to_bm25query(@description::text, 'issue_documents_description_bm25_idx')
    LIMIT 20
)
SELECT d.*,
       (COALESCE(1.0 / (60 + v.rank), 0.0) + COALESCE(1.0 / (60 + k.rank), 0.0))::double precision AS combined_score
FROM issue_documents d
LEFT JOIN vector_search v ON d.id = v.id
LEFT JOIN keyword_search k ON d.id = k.id
WHERE v.id IS NOT NULL OR k.id IS NOT NULL
ORDER BY combined_score DESC
LIMIT 10;
