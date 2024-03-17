-- +goose Up

CREATE TABLE IF NOT EXISTS transactions (
    id bigserial PRIMARY KEY,
    --
    billing_cycle_id bigserial NOT NULL,
    task_public_id uuid NOT NULL,
    account_public_id uuid NOT NULL,
    --
    credit NUMERIC(10,5) DEFAULT 0 NOT NULL,
    debit NUMERIC(10,5) DEFAULT 0 NOT NULL,
    --
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

-- +goose Down