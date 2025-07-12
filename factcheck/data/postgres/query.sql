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

-- name: UpdateTopic :one
UPDATE topics SET 
    name = $2,
    description = $3,
    status = $4,
    result = $5,
    result_status = $6,
    updated_at = $7
WHERE id = $1 RETURNING *;

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