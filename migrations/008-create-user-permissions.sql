-- +migrate Up
CREATE TABLE user_permissions (
    user_id uuid NOT NULL,
    permission_id uuid NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CONSTRAINT user_permissions_pkey PRIMARY KEY (user_id, permission_id)
);

CREATE INDEX idx_user_permissions_fk_permission_uid ON user_permissions USING btree (permission_id);

ALTER TABLE user_permissions
    ADD CONSTRAINT user_permissions_permission_id_fkey FOREIGN KEY (permission_id) REFERENCES permissions (id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE user_permissions
    ADD CONSTRAINT user_permissions_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE;

-- +migrate Down
DROP TABLE user_permissions;

