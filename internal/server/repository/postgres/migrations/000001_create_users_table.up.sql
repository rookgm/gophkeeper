CREATE TABLE IF NOT EXISTS "users" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "login" varchar NOT NULL UNIQUE,
    "password" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);