-- +migrate Up
CREATE TABLE access_tokens (
    token uuid NOT NULL,
    valid_until timestamptz NULL,
    user_id uuid NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CONSTRAINT access_tokens_pkey PRIMARY KEY (token)
);

CREATE INDEX idx_access_tokens_fk_user_uid ON access_tokens USING btree (user_id);

ALTER TABLE access_tokens
    ADD CONSTRAINT access_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE;

-- +migrate Down
DROP TABLE access_tokens;

