-- +goose Up

CREATE TABLE IF NOT EXISTS transactions (
    id bigserial PRIMARY KEY,
    --
    task_public_id uuid NOT NULL,
    account_public_id uuid NOT NULL,
    --
    amount NUMERIC(10,5) DEFAULT 0 NOT NULL,
    --
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

-- +goose Down