-- +migrate Up
CREATE TABLE refresh_tokens (
                               token uuid NOT NULL,
                               user_id uuid NULL,
                               created_at timestamptz NOT NULL,
                               updated_at timestamptz NOT NULL,
                               CONSTRAINT refresh_tokens_pkey PRIMARY KEY (token)
);
CREATE INDEX idx_refresh_tokens_fk_user_uid ON refresh_tokens USING btree (user_id);

ALTER TABLE refresh_tokens ADD CONSTRAINT refresh_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE;

-- +migrate Down
DROP TABLE refresh_tokens;