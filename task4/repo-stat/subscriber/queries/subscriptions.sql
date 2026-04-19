-- name: CreateSubscription :one
INSERT INTO subscriptions (owner, repo)
VALUES ($1, $2)
RETURNING id, owner, repo, created_at;

-- name: ListSubscriptions :many
SELECT owner, repo
FROM subscriptions
ORDER BY id DESC;

-- name: DeleteSubscription :exec
DELETE FROM subscriptions
WHERE owner = $1 AND repo = $2;