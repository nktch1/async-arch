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
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() not null,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() not null
);

CREATE UNIQUE INDEX ON tasks (public_id);

-- +goose Down
