-- +goose Up
CREATE TABLE IF NOT EXISTS campaign (
    campaign_id UUID PRIMARY KEY,
    advertiser_id UUID NOT NULL REFERENCES advertiser(advertiser_id),
    name TEXT NOT NULL,
    budget_total BIGINT NOT NULL,
    budget_daily BIGINT NOT NULL,
    spend_total BIGINT NOT NULL DEFAULT 0,
    spend_today BIGINT NOT NULL DEFAULT 0,
    spend_day DATE NOT NULL DEFAULT CURRENT_DATE,
    status TEXT NOT NULL DEFAULT 'active',
    pause_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT campaign_status_check CHECK (status IN ('active', 'paused')),
    CONSTRAINT campaign_pause_reason_check
        CHECK (pause_reason IN ('budget_daily_exceeded', 'budget_total_exceeded', 'manual')),
    CONSTRAINT campaign_pause_reason_consistency 
        CHECK ((status = 'paused') = (pause_reason IS NOT NULL))
);

-- +goose Down
DROP TABLE IF EXISTS campaign;