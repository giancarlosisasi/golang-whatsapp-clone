-- name: CreateParticipants :copyfrom
INSERT INTO conversation_participants (
    conversation_id,
    user_id
) VALUES (
    $1, $2
);

-- name: GetParticipantByUserAndConversation :one
SELECT * FROM conversation_participants
WHERE user_id = $1 AND conversation_id = $2 AND is_active = true;

-- name: UpdateParticipantLastReadAt :exec
UPDATE conversation_participants
SET
    last_read_at = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1 AND conversation_id = $3;

-- name: CountUnreadMessages :one
SELECT COUNT(*)::INTEGER
FROM messages m
JOIN conversation_participants cp ON m.conversation_id = cp.conversation_id
WHERE cp.user_id = $1
    AND cp.conversation_id = $2
    AND cp.is_active = true
    AND m.created_at > cp.last_read_at
    AND m.sender_id != $1;