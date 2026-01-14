-- name: GetUserFromRefreshToken :one
SELECT users.* FROM users
INNER JOIN refresh_tokens ON refresh_tokens.user_id = users.id
WHERE refresh_tokens.token = $1
AND revoked_at IS NULL
AND expires_at > NOW();
