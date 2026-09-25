-- name: GetCampaignByID :one
SELECT campaign_id, advertiser_id, name, budget_total, budget_daily,
       spend_total, spend_today, spend_day, status, pause_reason, created_at, updated_at
FROM campaign
WHERE campaign_id = $1;

-- name: GetCampaignByIDForUpdate :one
SELECT campaign_id, advertiser_id, name, budget_total, budget_daily,
       spend_total, spend_today, spend_day, status, pause_reason, created_at, updated_at
FROM campaign
WHERE campaign_id = $1
FOR UPDATE;

-- name: CreateCampaign :one
INSERT INTO campaign (campaign_id, advertiser_id, name, budget_total, budget_daily, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING campaign_id, advertiser_id, name, budget_total, budget_daily,
          spend_total, spend_today, spend_day, status, pause_reason, created_at, updated_at;

-- name: SaveCampaign :one
UPDATE campaign
SET name = $2, budget_total = $3, budget_daily = $4, updated_at = $5
WHERE campaign_id = $1
RETURNING campaign_id, advertiser_id, name, budget_total, budget_daily,
          spend_total, spend_today, spend_day, status, pause_reason, created_at, updated_at;

-- name: PauseCampaign :one
UPDATE campaign
SET status = 'paused', pause_reason = $2, updated_at = $3
WHERE campaign_id = $1 AND status = 'active'
RETURNING campaign_id, advertiser_id, name, budget_total, budget_daily,
          spend_total, spend_today, spend_day, status, pause_reason, created_at, updated_at;

-- name: ResumeCampaign :one
UPDATE campaign
SET status = 'active', pause_reason = NULL, updated_at = $2
WHERE campaign_id = $1 AND status = 'paused'
RETURNING campaign_id, advertiser_id, name, budget_total, budget_daily,
          spend_total, spend_today, spend_day, status, pause_reason, created_at, updated_at;

-- name: DeleteCampaign :exec
DELETE FROM campaign WHERE campaign_id = $1;
