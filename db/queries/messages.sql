-- name: CreateMessage :one
INSERT INTO bot_history_messages (
  bot_id,
  route_id,
  sender_channel_identity_id,
  sender_account_user_id,
  channel_type,
  source_message_id,
  source_reply_to_message_id,
  role,
  content,
  metadata
)
VALUES (
  sqlc.arg(bot_id),
  sqlc.narg(route_id)::uuid,
  sqlc.narg(sender_channel_identity_id)::uuid,
  sqlc.narg(sender_user_id)::uuid,
  sqlc.narg(platform)::text,
  sqlc.narg(external_message_id)::text,
  sqlc.narg(source_reply_to_message_id)::text,
  sqlc.arg(role),
  sqlc.arg(content),
  sqlc.arg(metadata)
)
RETURNING
  id,
  bot_id,
  route_id,
  sender_channel_identity_id,
  sender_account_user_id AS sender_user_id,
  channel_type AS platform,
  source_message_id AS external_message_id,
  source_reply_to_message_id,
  role,
  content,
  metadata,
  created_at;

-- name: ListMessages :many
SELECT
  id,
  bot_id,
  route_id,
  sender_channel_identity_id,
  sender_account_user_id AS sender_user_id,
  channel_type AS platform,
  source_message_id AS external_message_id,
  source_reply_to_message_id,
  role,
  content,
  metadata,
  created_at
FROM bot_history_messages
WHERE bot_id = sqlc.arg(bot_id)
ORDER BY created_at ASC;

-- name: ListMessagesSince :many
SELECT
  id,
  bot_id,
  route_id,
  sender_channel_identity_id,
  sender_account_user_id AS sender_user_id,
  channel_type AS platform,
  source_message_id AS external_message_id,
  source_reply_to_message_id,
  role,
  content,
  metadata,
  created_at
FROM bot_history_messages
WHERE bot_id = sqlc.arg(bot_id)
  AND created_at >= sqlc.arg(created_at)
ORDER BY created_at ASC;

-- name: ListMessagesBefore :many
SELECT
  id,
  bot_id,
  route_id,
  sender_channel_identity_id,
  sender_account_user_id AS sender_user_id,
  channel_type AS platform,
  source_message_id AS external_message_id,
  source_reply_to_message_id,
  role,
  content,
  metadata,
  created_at
FROM bot_history_messages
WHERE bot_id = sqlc.arg(bot_id)
  AND created_at < sqlc.arg(created_at)
ORDER BY created_at DESC
LIMIT sqlc.arg(max_count);

-- name: ListMessagesLatest :many
SELECT
  id,
  bot_id,
  route_id,
  sender_channel_identity_id,
  sender_account_user_id AS sender_user_id,
  channel_type AS platform,
  source_message_id AS external_message_id,
  source_reply_to_message_id,
  role,
  content,
  metadata,
  created_at
FROM bot_history_messages
WHERE bot_id = sqlc.arg(bot_id)
ORDER BY created_at DESC
LIMIT sqlc.arg(max_count);

-- name: DeleteMessagesByBot :exec
DELETE FROM bot_history_messages
WHERE bot_id = sqlc.arg(bot_id);
