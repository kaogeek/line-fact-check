-- name: CreateTopic :one
INSERT INTO topics (
    id, name, description, status, result, result_status, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetTopic :one
SELECT * FROM topics WHERE id = $1;

-- name: ListTopics :many
WITH numbered_topics AS (
    SELECT *, 
           ROW_NUMBER() OVER (ORDER BY created_at DESC) as rn,
           COUNT(*) OVER () as total_count
    FROM topics
)
SELECT id, name, description, status, result, result_status, created_at, updated_at
FROM numbered_topics
WHERE CASE 
    WHEN $1 = 0 THEN true  -- No pagination
    WHEN $1 > 0 THEN rn BETWEEN $2 + 1 AND $2 + $1  -- Normal pagination
    WHEN $1 < 0 THEN rn BETWEEN total_count + $1 + 1 AND total_count + $2  -- Negative pagination
END
ORDER BY created_at DESC;

-- name: ListTopicsByStatus :many
WITH numbered_topics AS (
    SELECT *, 
           ROW_NUMBER() OVER (ORDER BY created_at DESC) as rn,
           COUNT(*) OVER () as total_count
    FROM topics
    WHERE status = $1
)
SELECT id, name, description, status, result, result_status, created_at, updated_at
FROM numbered_topics
WHERE CASE 
    WHEN $2 = 0 THEN true  -- No pagination
    WHEN $2 > 0 THEN rn BETWEEN $3 + 1 AND $3 + $2  -- Normal pagination
    WHEN $2 < 0 THEN rn BETWEEN total_count + $2 + 1 AND total_count + $3  -- Negative pagination
END
ORDER BY created_at DESC;

-- name: ListTopicsInIDs :many
SELECT DISTINCT t.* FROM topics t 
WHERE t.id = ANY($1::uuid[]) 
ORDER BY t.created_at DESC;

-- name: ListTopicsLikeID :many
WITH numbered_topics AS (
    SELECT *, 
           ROW_NUMBER() OVER (ORDER BY created_at DESC) as rn,
           COUNT(*) OVER () as total_count
    FROM topics t 
    WHERE t.id::text LIKE $1::text
)
SELECT id, name, description, status, result, result_status, created_at, updated_at
FROM numbered_topics
WHERE CASE 
    WHEN $2 = 0 THEN true  -- No pagination
    WHEN $2 > 0 THEN rn BETWEEN $3 + 1 AND $3 + $2  -- Normal pagination
    WHEN $2 < 0 THEN rn BETWEEN total_count + $2 + 1 AND total_count + $3  -- Negative pagination
END
ORDER BY created_at DESC;

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

-- name: CountTopicsGroupByStatusLikeID :many
SELECT status, COUNT(*) as count 
FROM topics t 
WHERE t.id::text LIKE $1::text 
GROUP BY status;

-- name: CountTopicsGroupByStatusLikeMessageText :many
SELECT t.status, COUNT(DISTINCT t.id) as count 
FROM topics t 
INNER JOIN messages m ON t.id = m.topic_id 
WHERE m.text ILIKE $1 
GROUP BY t.status;

-- name: CountTopicsGroupByStatusLikeIDLikeMessageText :many
SELECT t.status, COUNT(DISTINCT t.id) as count 
FROM topics t 
INNER JOIN messages m ON t.id = m.topic_id 
WHERE t.id::text LIKE $1::text AND m.text ILIKE $2 
GROUP BY t.status;

-- name: DeleteTopic :exec
DELETE FROM topics WHERE id = $1;

-- name: CreateMessage :one
INSERT INTO messages (
    id, user_message_id, type, status, topic_id, text, language, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetMessage :one
SELECT * FROM messages WHERE id = $1;

-- name: ListMessagesByTopic :many
SELECT * FROM messages WHERE topic_id = $1 ORDER BY created_at ASC;

-- name: AssignMessageToTopic :one
UPDATE messages SET 
    topic_id = $2,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: DeleteMessage :exec
DELETE FROM messages WHERE id = $1;

-- name: CreateUserMessage :one
INSERT INTO user_messages (
    id, type, replied_at, metadata, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetUserMessage :one
SELECT * FROM user_messages WHERE id = $1;

-- name: DeleteUserMessage :exec
DELETE FROM user_messages WHERE id = $1; 

-- name: ListTopicsDynamic :many
SELECT DISTINCT t.*
FROM topics t
LEFT JOIN messages m ON t.id = m.topic_id
WHERE 1=1
    AND CASE 
        WHEN $1::text != '' THEN t.id::text LIKE $1::text 
        ELSE true 
    END
    AND CASE 
        WHEN array_length($2::text[], 1) > 0 THEN t.status = ANY($2::text[])
        ELSE true 
    END
    AND CASE 
        WHEN $3::text != '' THEN (
            CASE 
                WHEN m.language = 'th' THEN m.text LIKE $3::text COLLATE "C"
                WHEN m.language = 'en' THEN m.text ILIKE $3::text
                ELSE m.text ILIKE $3::text  -- fallback for unknown language
            END
        )
        ELSE true 
    END
ORDER BY t.created_at DESC
LIMIT CASE WHEN $4::integer = 0 THEN NULL ELSE $4::integer END
OFFSET CASE WHEN $4::integer = 0 THEN 0 ELSE $5::integer END;

-- name: CountTopicsGroupByStatusDynamic :many
SELECT t.status, COUNT(DISTINCT t.id) as count 
FROM topics t 
LEFT JOIN messages m ON t.id = m.topic_id 
WHERE 1=1
    AND CASE 
        WHEN $1::text != '' THEN t.id::text LIKE $1::text 
        ELSE true 
    END
    AND CASE 
        WHEN $2::text != '' THEN (
            CASE 
                WHEN m.language = 'th' THEN m.text LIKE $2::text COLLATE "C"
                WHEN m.language = 'en' THEN m.text ILIKE $2::text
                ELSE m.text ILIKE $2::text  -- fallback for unknown language
            END
        )
        ELSE true 
    END
GROUP BY t.status;

-- name: CreateMessageV2 :one
INSERT INTO messages_v2 (
    id, user_id, topic_id, type, text, language, metadata, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetMessageV2 :one
SELECT * FROM messages_v2 WHERE id = $1;

-- name: ListMessagesV2ByTopic :many
SELECT * FROM messages_v2 WHERE topic_id = $1 ORDER BY created_at ASC;

-- name: ListMessagesV2ByGroup :many
SELECT * FROM messages_v2 WHERE group_id = $1 ORDER BY created_at ASC;

-- name: AssignMessageV2ToTopic :one
UPDATE messages_v2 SET 
    topic_id = $2,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: UnassignMessageV2FromTopic :one
UPDATE messages_v2 SET 
    topic_id = NULL,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: DeleteMessageV2 :exec
DELETE FROM messages_v2 WHERE id = $1; 


-- name: CreateMessageV2Group :one
INSERT INTO messages_v2_group (
    id, topic_id, name, text, text_sha1, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetMessageV2Group :one
SELECT * FROM messages_v2_group WHERE id = $1;

-- name: ListMessageV2GroupsByTopic :many
SELECT * FROM messages_v2_group WHERE topic_id = $1 ORDER BY created_at ASC;

-- name: UpdateMessageV2GroupName :one
UPDATE messages_v2_group SET 
    name = $2,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: AssignMessageV2GroupToTopic :one
UPDATE messages_v2_group SET 
    topic_id = $2,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: UnassignMessageV2GroupFromTopic :one
UPDATE messages_v2_group SET 
    topic_id = NULL,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: DeleteMessageV2Group :exec
DELETE FROM messages_v2_group WHERE id = $1;

-- name: AssignMessageV2Group :one
UPDATE messages_v2 SET 
    group_id = $2,
    updated_at = NOW()
WHERE id = $1 RETURNING *;
