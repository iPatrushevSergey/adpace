-- name: GetByID :one
SELECT advertiser_id, name, country, created_at, updated_at
FROM advertiser
WHERE advertiser_id = $1;

-- name: Create :one
INSERT INTO advertiser (advertiser_id, name, country, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING advertiser_id, name, country, created_at, updated_at;

-- name: Put :one
UPDATE advertiser
SET name = $2, country = $3, updated_at = $4
WHERE advertiser_id = $1
RETURNING advertiser_id, name, country, created_at, updated_at;

-- name: Delete :exec
DELETE FROM advertiser WHERE advertiser_id = $1;