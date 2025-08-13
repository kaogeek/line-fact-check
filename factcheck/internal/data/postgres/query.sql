-- name: CreateTopic :one
INSERT INTO topics (
    id, name, description, status, result, result_status, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetTopic :one
SELECT * FROM topics WHERE id = $1;

-- name: GetTopicStatus :one
SELECT status FROM topics WHERE id = $1;

-- name: TopicExists :one
SELECT EXISTS (SELECT 1 from topics where id = $1);

-- name: ListTopics :many
SELECT DISTINCT t.*
FROM topics t
ORDER BY t.created_at DESC
LIMIT CASE WHEN $1::integer = 0 THEN NULL ELSE $1::integer END
OFFSET CASE WHEN $2::integer = 0 THEN 0 ELSE $2::integer END;

-- name: ListTopicsByStatus :many
SELECT DISTINCT t.*
FROM topics t
WHERE status = $1
ORDER BY t.created_at DESC
LIMIT CASE WHEN $2::integer = 0 THEN NULL ELSE $2::integer END
OFFSET CASE WHEN $3::integer = 0 THEN 0 ELSE $3::integer END;

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

-- name: ResolveTopic :one
UPDATE topics SET
    result = $2,
    result_status = $3,
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

-- name: ListTopicsDynamicV2 :many
SELECT DISTINCT t.*
FROM topics t
LEFT JOIN message_groups mg ON t.id = mg.topic_id
LEFT JOIN messages_v2 m ON mg.id = m.group_id
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
                WHEN mg.language = 'th' THEN m.text LIKE $3::text COLLATE "C"
                WHEN mg.language = 'en' THEN m.text ILIKE $3::text
                ELSE m.text ILIKE $3::text  -- fallback for unknown language
            END
        )
        ELSE true
    END
ORDER BY t.created_at DESC
LIMIT CASE WHEN $4::integer = 0 THEN NULL ELSE $4::integer END
OFFSET CASE WHEN $4::integer = 0 THEN 0 ELSE $5::integer END;

-- name: CountTopicsGroupByStatusDynamicV2 :many
SELECT t.status, COUNT(DISTINCT t.id) as count
FROM topics t
LEFT JOIN message_groups m ON t.id = m.topic_id
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
    id, user_id, topic_id, group_id, type_user, type, text, language, metadata, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
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

-- name: AssignMessageV2ToMessageGroup :one
UPDATE messages_v2 SET
    group_id = $2,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: DeleteMessageV2 :exec
DELETE FROM messages_v2 WHERE id = $1;

-- name: CreateMessageGroup :one
INSERT INTO message_groups (
    id, status, topic_id, name, text, text_sha1, language, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: ListMessageGroupDynamic :many
SELECT  mg.*
FROM message_groups mg
WHERE 1=1
    AND CASE
        WHEN sqlc.arg('text')::text != '' THEN mg.text::text LIKE sqlc.arg('text')::text
        ELSE true
    END
    AND CASE
        WHEN array_length(sqlc.arg('id_in')::text[], 1) > 0 THEN mg.id = ANY((sqlc.arg('id_in')::text[])::uuid[])
        ELSE true
    END
    AND CASE
        WHEN array_length(sqlc.arg('id_not_in')::text[], 1) > 0 THEN NOT (mg.id = ANY((sqlc.arg('id_not_in')::text[])::uuid[]))
        ELSE true
    END
    AND CASE
        WHEN array_length(sqlc.arg('status')::text[], 1) > 0 THEN mg.status = ANY(sqlc.arg('status')::text[])
        ELSE true
    END
ORDER BY mg.created_at DESC
LIMIT CASE WHEN sqlc.arg('limit')::integer = 0 THEN NULL ELSE sqlc.arg('limit')::integer END
OFFSET CASE WHEN sqlc.arg('offset')::integer = 0 THEN 0 ELSE sqlc.arg('offset')::integer END;

-- name: GetMessageGroup :one
SELECT * FROM message_groups WHERE id = $1;

-- name: GetMessageGroupBySHA1 :one
SELECT * FROM message_groups WHERE text_sha1 = $1;

-- name: ListMessageGroupsByTopic :many
SELECT * FROM message_groups WHERE topic_id = $1 ORDER BY created_at ASC;

-- name: UpdateMessageGroupName :one
UPDATE message_groups SET
    name = $2,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: UpdateMessageGroupStatus :one
UPDATE message_groups SET
    status = $2,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: AssignMessageGroupToTopic :one
UPDATE message_groups SET
    topic_id = $2,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: UnassignMessageGroupFromTopic :one
UPDATE message_groups SET
    topic_id = NULL,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: DeleteMessageGroup :exec
DELETE FROM message_groups WHERE id = $1;

-- name: CreateAnswer :one
INSERT INTO answers (
    id, topic_id, text, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetAnswerByID :one
SELECT * FROM answers WHERE id = $1;

-- name: GetAnswerByTopicID :one
SELECT * FROM answers WHERE topic_id = $1 ORDER BY created_at DESC LIMIT 1;

-- name: ListAnswersByTopicID :many
SELECT * FROM answers WHERE topic_id = $1 ORDER BY created_at DESC;

-- name: DeleteAnswer :exec
DELETE FROM answers WHERE id = $1;
