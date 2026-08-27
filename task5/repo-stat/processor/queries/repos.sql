-- name: GetInfoByOwnerAndRepo :one
SELECT id, owner, repo, full_name, description, stars, forks, created_at, updated_at
FROM repositories
WHERE owner = $1 AND repo = $2;

-- name: UpsertInfo :exec
INSERT INTO repositories (
    owner, repo, full_name, description, stars, forks, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, NOW()
    )
ON CONFLICT (owner, repo)
DO UPDATE SET
    full_name = EXCLUDED.full_name,
    description = EXCLUDED.description,
    stars = EXCLUDED.stars,
    forks = EXCLUDED.forks,
    updated_at = NOW();

-- name: ListAllStats :many
SELECT id, owner, repo, full_name, description, stars, forks, created_at, updated_at
FROM repositories
ORDER BY updated_at DESC;