-- +goose Up

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE task_status AS ENUM ('new', 'done');

CREATE TABLE IF NOT EXISTS tasks (
    id bigserial PRIMARY KEY,
    public_id uuid DEFAULT gen_random_uuid() NOT NULL,
    account_public_id uuid NOT NULL,
    --
    description VARCHAR NOT NULL,
    --
    status task_status NOT NULL,
    --
    cost NUMERIC(10,5) DEFAULT 0 NOT NULL,
    --
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

CREATE UNIQUE INDEX ON tasks (public_id);

-- +goose Down
