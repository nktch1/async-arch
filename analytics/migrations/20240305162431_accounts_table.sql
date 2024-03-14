-- +goose Up

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE "account_role" AS ENUM (
    'employee', 'manager', 'accountant', 'admin'
);

CREATE TABLE IF NOT EXISTS "accounts" (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID DEFAULT gen_random_uuid() NOT NULL,
    --
    name VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    --
    role account_role NOT NULL,
    --
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

CREATE UNIQUE INDEX ON "accounts" (public_id);
CREATE UNIQUE INDEX ON "accounts" (email);

-- +goose Down
