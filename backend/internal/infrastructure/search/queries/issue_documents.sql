-- name: CreateIssueDocument :one
INSERT INTO issue_documents (id, workspace_id, title, description, created_at, updated_at)
VALUES (@id, @workspace_id, @title, @description, NOW(), NOW())
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
RETURNING *;

-- name: FindIssueDocumentsByDescription :many
SELECT * FROM issue_documents
ORDER BY description <@> to_bm25query(@description::text, 'issue_documents_description_bm25_idx')
LIMIT 50;
