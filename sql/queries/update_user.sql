-- name: UpdateUser :one
UPDATE users
SET hashed_password = $2, email = $3
WHERE id = $1
RETURNING *;