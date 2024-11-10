-- name: CreateFeedFollow :many
with inserted AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
) 
SELECT a.*, 
b.name AS feed_name, 
c.name AS user_name
FROM inserted a 
INNER JOIN feeds b ON a.feed_id = b.id 
INNER JOIN users c ON a.user_id = c.id;

-- name: GetFeedFollowsForUser :many
SELECT a.*,
b.name as feed_name,
c.name as user_name
FROM feed_follows a 
INNER JOIN feeds b ON a.feed_id = b.id
INNER JOIN users c ON a.user_id = c.id
WHERE a.user_id = $1;

-- name: DeleteFeedFollow :one
DELETE FROM feed_follows 
WHERE feed_id IN (
    SELECT feed_follows.feed_id FROM  feed_follows
    JOIN feeds ON feed_follows.feed_id = feeds.id
    WHERE feeds.url = $1
) AND feed_follows.user_id = $2
RETURNING *;
