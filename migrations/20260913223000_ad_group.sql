-- +goose Up
CREATE TABLE IF NOT EXISTS ad_group (
    ad_group_id UUID PRIMARY KEY,
    campaign_id UUID NOT NULL REFERENCES campaign(campaign_id),
    name TEXT NOT NULL,
    target_country CHAR(2),
    status TEXT NOT NULL DEFAULT 'active', -- permission to display
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT ad_group_status_check CHECK (status IN ('active', 'paused'))
);

-- +goose Down
DROP TABLE IF EXISTS ad_group;