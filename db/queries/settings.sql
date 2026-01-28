-- name: GetSettingsByUserID :one
SELECT user_id, max_context_load_time, language
FROM user_settings
WHERE user_id = $1;

-- name: UpsertSettings :one
INSERT INTO user_settings (user_id, max_context_load_time, language)
VALUES ($1, $2, $3)
ON CONFLICT (user_id) DO UPDATE SET
  max_context_load_time = EXCLUDED.max_context_load_time,
  language = EXCLUDED.language
RETURNING user_id, max_context_load_time, language;

-- name: DeleteSettingsByUserID :exec
DELETE FROM user_settings
WHERE user_id = $1;

