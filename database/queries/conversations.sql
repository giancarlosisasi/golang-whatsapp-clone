-- name: CreateConversation :one
INSERT INTO conversations (type)
VALUES ($1)
RETURNING *;

-- name: GetUserConversations :many
SELECT
    c.id,
    c.type,
    c.last_message_at,
    c.created_at,
    c.updated_at,
    (
        SELECT COUNT(*)::INTEGER
        FROM messages unread_m
        WHERE unread_m.conversation_id = c.id
            AND unread_m.created_at > cp.last_read_at
            AND unread_m.sender_id != $1
    ) as unread_count
FROM conversations c
JOIN conversation_participants cp ON c.id = cp.conversation_id
WHERE cp.user_id = $1 AND cp.is_active = true
ORDER BY c.last_message_at DESC NULLS LAST;

-- -- name: GetConversationParticipants :many
-- SELECT
--     cp.id,
--     cp.conversation_id,
--     cp.user_id,
--     cp.joined_at,
--     cp.last_read_at,
--     cp.is_active,
--     u.id as user_id,
--     u.name as user_name,
--     u.avatar_url as user_avatar_url
-- FROM conversation_participants cp
-- JOIN users u ON cp.user_id = u.id
-- WHERE cp.conversation_id = ANY($1::uuid[])
--     AND cp.is_active = true
-- ORDER BY cp.joined_at ASC;

-- name: GetLastMessage :many
WITH ranked_messages AS (
    SELECT
        m.id,
        m.conversation_id,
        m.sender_id,
        m.content,
        m.message_type,
        m.created_at,
        m.edited_at,
        m.reply_to_message_id,
        sender.name as sender_name,
        sender.avatar_url as sender_avatar_url,
        sender.email as sender_email,
        sender.created_at as sender_created_at,
        sender.updated_at as sender_updated_at,
        ROW_NUMBER() OVER (PARTITION BY m.conversation_id ORDER BY m.created_at DESC) as rn
    FROM messages m
    JOIN users sender ON m.sender_id = sender.id
    WHERE m.conversation_id = ANY($1::uuid[])
)
SELECT
    id,
    conversation_id,
    sender_id,
    content,
    message_type,
    created_at,
    edited_at,
    reply_to_message_id,
    sender_name,
    sender_avatar_url,
    sender_email,
    sender_created_at,
    sender_updated_at
FROM ranked_messages
WHERE rn = 1
ORDER BY created_at DESC;

-- name: FindDirectConversation :one
SELECT c.*
FROM conversations c
WHERE c.type = 'DIRECT'
    AND EXISTS (
        SELECT 1 FROM conversation_participants cp1
        WHERE cp1.conversation_id = c.id AND cp1.user_id = $1 AND cp1.is_active = true
    )
    AND EXISTS (
        SELECT 1 FROM conversation_participants cp2
        WHERE cp2.conversation_id = c.id AND cp2.user_id = $2 AND cp2.is_active = true
    )
    AND (
        SELECT COUNT(*) FROM conversation_participants cp
        WHERE cp.conversation_id = c.id AND cp.is_active = true
    ) = 2;

-- name: UpdateConversationLastMessageAt :exec
UPDATE conversations
SET
    last_message_at = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;
