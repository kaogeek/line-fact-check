-- name: CreateTopic :one
INSERT INTO topics (
    id, name, description, status, result, result_status, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetTopic :one
SELECT * FROM topics WHERE id = $1;

-- name: ListTopics :many
SELECT * FROM topics ORDER BY created_at DESC;

-- name: ListTopicsByStatus :many
SELECT * FROM topics WHERE status = $1 ORDER BY created_at DESC;

-- name: ListTopicsInIDs :many
SELECT DISTINCT t.* FROM topics t 
WHERE t.id = ANY($1::uuid[]) 
ORDER BY t.created_at DESC;

-- name: ListTopicsByMessageText :many
SELECT DISTINCT t.* FROM topics t 
INNER JOIN messages m ON t.id = m.topic_id 
WHERE m.text LIKE $1 
ORDER BY t.created_at DESC;

-- name: ListTopicsInIDsAndMessageText :many
SELECT DISTINCT t.* FROM topics t 
INNER JOIN messages m ON t.id = m.topic_id 
WHERE t.id = ANY($1::uuid[]) AND m.text LIKE $2 
ORDER BY t.created_at DESC;

-- name: ListTopicsLikeID :many
SELECT * FROM topics t 
WHERE t.id::text LIKE $1::text 
ORDER BY t.created_at DESC;

-- name: ListTopicsLikeIDAndMessageText :many
SELECT DISTINCT t.* FROM topics t 
INNER JOIN messages m ON t.id = m.topic_id 
WHERE t.id::text LIKE $1::text AND m.text LIKE $2 
ORDER BY t.created_at DESC;

-- name: ListTopicsByStatusAndMessageText :many
SELECT DISTINCT t.* FROM topics t 
INNER JOIN messages m ON t.id = m.topic_id 
WHERE t.status = $1 AND m.text LIKE $2 
ORDER BY t.created_at DESC;

-- name: ListTopicsByStatusAndLikeID :many
SELECT * FROM topics t 
WHERE t.status = $1 AND t.id::text LIKE $2::text 
ORDER BY t.created_at DESC;

-- name: ListTopicsByStatusAndLikeIDAndMessageText :many
SELECT DISTINCT t.* FROM topics t 
INNER JOIN messages m ON t.id = m.topic_id 
WHERE t.status = $1 AND t.id::text LIKE $2::text AND m.text LIKE $3 
ORDER BY t.created_at DESC;

-- name: UpdateTopicStatus :one
UPDATE topics SET 
    status = $2,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: UpdateTopicDescription :one
UPDATE topics SET 
    description = $2,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: UpdateTopicName :one
UPDATE topics SET 
    name = $2,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: CountTopicsByStatus :one
SELECT COUNT(*) FROM topics WHERE status = $1;

-- name: CountTopicsGroupedByStatus :many
SELECT status, COUNT(*) as count 
FROM topics 
GROUP BY status;

-- name: DeleteTopic :exec
DELETE FROM topics WHERE id = $1;

-- name: CreateMessage :one
INSERT INTO messages (
    id, topic_id, text, type, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetMessage :one
SELECT * FROM messages WHERE id = $1;

-- name: ListMessagesByTopic :many
SELECT * FROM messages WHERE topic_id = $1 ORDER BY created_at ASC;

-- name: UpdateMessage :one
UPDATE messages SET 
    text = $2,
    type = $3,
    updated_at = $4
WHERE id = $1 RETURNING *;

-- name: AssignMessageToTopic :one
UPDATE messages SET 
    topic_id = $2,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: DeleteMessage :exec
DELETE FROM messages WHERE id = $1;

-- name: CreateUserMessage :one
INSERT INTO user_messages (
    id, replied_at, message_id, metadata, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetUserMessage :one
SELECT * FROM user_messages WHERE id = $1;

-- name: ListUserMessagesByMessage :many
SELECT * FROM user_messages WHERE message_id = $1 ORDER BY created_at ASC;

-- name: UpdateUserMessage :one
UPDATE user_messages SET 
    replied_at = $2,
    metadata = $3,
    updated_at = $4
WHERE id = $1 RETURNING *;

-- name: DeleteUserMessage :exec
DELETE FROM user_messages WHERE id = $1; 