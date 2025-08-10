-- name: CreateConversation :one
INSERT INTO conversations (type)
VALUES ($1)
RETURNING *;

-- name: GetUserConversations :many
SELECT
    c.*,
    m.content as last_message_content,
    m.created_at as last_message_created_at,
    sender.name as last_message_sender_name
FROM conversations c
JOIN conversation_participants cp ON c.id = cp.conversation_id
LEFT JOIN messages m ON c.id = m.conversation_id
    AND m.created_at = c.last_message_at
LEFT JOIN users sender ON m.sender_id = sender.id
WHERE cp.user_id = $1 AND cp.is_active = true
ORDER BY c.last_message_at DESC;

-- name: FindDirectConversation :one
SELECT c.*
FROM conversations c
WHERE c.type = 'direct'
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
