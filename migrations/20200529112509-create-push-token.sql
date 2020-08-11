-- +migrate Up
CREATE TYPE provider_type AS ENUM (
    'fcm',
    'apn'
);

CREATE TABLE push_tokens (
    id uuid NOT NULL DEFAULT uuid_generate_v4 (),
    token text NOT NULL,
    provider provider_type NOT NULL,
    user_id uuid NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CONSTRAINT push_tokens_pkey PRIMARY KEY (id),
    CONSTRAINT push_tokens_token_key UNIQUE (token)
);

CREATE INDEX idx_push_tokens_fk_user_id ON push_tokens USING btree (user_id);

CREATE INDEX idx_push_tokens_token ON push_tokens USING btree (token);

ALTER TABLE push_tokens
    ADD CONSTRAINT push_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE;

-- +migrate Down
DROP TABLE IF EXISTS push_tokens;

DROP TYPE IF EXISTS provider_type;

