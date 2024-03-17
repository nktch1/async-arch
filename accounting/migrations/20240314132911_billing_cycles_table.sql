-- +goose Up

CREATE TYPE billing_cycle_status AS ENUM ('new', 'done');

CREATE TABLE IF NOT EXISTS billing_cycles (
    id bigserial PRIMARY KEY,
    --
    start_date TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE,
    status billing_cycle_status NOT NULL,
    --
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

-- +goose Down