-- name: CreateMessage :one
INSERT INTO messages (
    conversation_id,
    sender_id,
    content,
    message_type,
    reply_to_message_id
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetConversationMessages :many
SELECT
    m.*,
    sender.name as sender_name,
    sender.avatar_url as sender_avatar_url,
    reply_msg.content as reply_content,
    reply_sender.name as reply_sender_name
FROM messages m
JOIN users sender ON m.sender_id = sender.id
LEFT JOIN messages reply_msg ON m.reply_to_message_id = reply_msg.id
LEFT JOIN users reply_sender ON reply_msg.sender_id = reply_sender.id
WHERE m.conversation_id = $1
ORDER BY m.created_at DESC
LIMIT $2 OFFSET $3;