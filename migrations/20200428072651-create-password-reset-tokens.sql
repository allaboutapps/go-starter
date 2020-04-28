-- +migrate Up
CREATE TABLE password_reset_tokens (
    token uuid NOT NULL DEFAULT uuid_generate_v4 (),
    valid_until timestamptz NOT NULL,
    user_id uuid NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CONSTRAINT password_reset_tokens_pkey PRIMARY KEY (token)
);

CREATE INDEX idx_password_reset_tokens_fk_user_uid ON password_reset_tokens USING btree (user_id);

ALTER TABLE password_reset_tokens
    ADD CONSTRAINT password_reset_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE;

-- +migrate Down
DROP TABLE IF EXISTS password_reset_tokens;

