-- name: InsertIssueLabel :exec
INSERT INTO issue_labels (issue_id, label_id) VALUES (@issue_id, @label_id);

-- name: ListLabelsByIssueIDs :many
SELECT il.issue_id, l.id, l.name
FROM issue_labels il
JOIN labels l ON l.id = il.label_id
WHERE il.issue_id = ANY(@issue_ids::uuid[])
ORDER BY l.name ASC;

-- name: GetLabelsByIDs :many
SELECT id, name FROM labels
WHERE id = ANY(@ids::uuid[]);
