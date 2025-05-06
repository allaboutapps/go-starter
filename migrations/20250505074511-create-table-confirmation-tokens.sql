-- +migrate Up
CREATE TABLE confirmation_tokens (
    token uuid NOT NULL DEFAULT uuid_generate_v4 (),
    valid_until timestamptz NOT NULL,
    user_id uuid NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CONSTRAINT confirmation_tokens_pkey PRIMARY KEY (token)
);

CREATE INDEX idx_confirmation_tokens_fk_user_uid ON confirmation_tokens USING btree (user_id);

ALTER TABLE confirmation_tokens
    ADD CONSTRAINT confirmation_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE users
    ADD COLUMN requires_confirmation boolean NOT NULL DEFAULT FALSE;

-- +migrate Down
ALTER TABLE users
    DROP COLUMN requires_confirmation;

DROP TABLE IF EXISTS confirmation_tokens;

