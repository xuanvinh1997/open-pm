-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL;

-- name: CreateUser :one
INSERT INTO users (email, encrypted_password, first_name, last_name, display_name, email_confirmed_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateUser :one
UPDATE users SET
    first_name = COALESCE(sqlc.narg('first_name'), first_name),
    last_name = COALESCE(sqlc.narg('last_name'), last_name),
    display_name = COALESCE(sqlc.narg('display_name'), display_name),
    avatar_url = COALESCE(sqlc.narg('avatar_url'), avatar_url)
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users SET encrypted_password = $2 WHERE id = $1;

-- name: SetUserEmailConfirmed :exec
UPDATE users SET email_confirmed_at = NOW() WHERE id = $1;

-- name: SetUserLastSignIn :exec
UPDATE users SET last_sign_in_at = NOW() WHERE id = $1;

-- name: SetUserConfirmationToken :exec
UPDATE users SET confirmation_token = $2, confirmation_sent_at = NOW() WHERE id = $1;

-- name: GetUserByConfirmationToken :one
SELECT * FROM users WHERE confirmation_token = $1 AND deleted_at IS NULL;

-- name: SetUserRecoveryToken :exec
UPDATE users SET recovery_token = $2, recovery_sent_at = NOW() WHERE id = $1;

-- name: GetUserByRecoveryToken :one
SELECT * FROM users WHERE recovery_token = $1 AND deleted_at IS NULL;

-- name: CreateSession :one
INSERT INTO sessions (user_id, ip, user_agent)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetSessionByID :one
SELECT * FROM sessions WHERE id = $1;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE id = $1;

-- name: DeleteUserSessions :exec
DELETE FROM sessions WHERE user_id = $1;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (session_id, user_id, token, parent)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetRefreshTokenByToken :one
SELECT * FROM refresh_tokens WHERE token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens SET revoked = TRUE WHERE token = $1;

-- name: RevokeRefreshTokenFamily :exec
UPDATE refresh_tokens SET revoked = TRUE WHERE session_id = $1;

-- name: CreateIdentity :one
INSERT INTO identities (user_id, provider, provider_id, identity_data)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetIdentityByProvider :one
SELECT * FROM identities WHERE provider = $1 AND provider_id = $2;

-- name: CreateOneTimeToken :one
INSERT INTO one_time_tokens (user_id, token_hash, token_type, relates_to)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetOneTimeTokenByHash :one
SELECT * FROM one_time_tokens WHERE token_hash = $1 AND token_type = $2;

-- name: DeleteOneTimeToken :exec
DELETE FROM one_time_tokens WHERE id = $1;

-- name: DeleteOneTimeTokensByUser :exec
DELETE FROM one_time_tokens WHERE user_id = $1 AND token_type = $2;
