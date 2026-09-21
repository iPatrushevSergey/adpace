-- +goose Up
CREATE TABLE IF NOT EXISTS advertiser (
    advertiser_id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    country CHAR(2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS advertiser;