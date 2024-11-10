-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeed :many
SELECT a.name, a.url, b.name AS username 
FROM feeds a
LEFT JOIN users b ON a.user_id = b.id;

-- name: GetFeedIDByUrl :one
SELECT id 
FROM feeds
WHERE url = $1 LIMIT 1;
