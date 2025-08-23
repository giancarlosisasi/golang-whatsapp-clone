-- name: CreateMessage :one
INSERT INTO messages (
    conversation_id,
    sender_id,
    content,
    message_type,
    reply_to_message_id,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetConversationMessages :many
SELECT
    m.*,
    sender.id as sender_id,
    sender.name as sender_name,
    sender.email as sender_email,
    sender.avatar_url as sender_avatar_url,
    sender.created_at as sender_created_at,
    sender.updated_at as sender_updated_at,
    reply_msg.id as reply_id,
    reply_msg.content as reply_content,
    reply_msg.message_type as reply_message_type,
    reply_sender.name as reply_sender_name
FROM messages m
JOIN users sender ON m.sender_id = sender.id
LEFT JOIN messages reply_msg ON m.reply_to_message_id = reply_msg.id
LEFT JOIN users reply_sender ON reply_msg.sender_id = reply_sender.id
WHERE m.conversation_id = $1
ORDER BY m.created_at DESC
LIMIT $2 OFFSET $3;