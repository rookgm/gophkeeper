CREATE TABLE IF NOT EXISTS "secrets" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    "name" varchar NOT NULL,
    "type" integer NOT NULL,
    "note" varchar DEFAULT '',
    "data" bytea NOT NULL,
    "deleted" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);