-- +migrate Up
CREATE TABLE user_roles (
                                 user_id uuid NOT NULL,
                                 role_id uuid NOT NULL,
                                 created_at timestamptz NOT NULL,
                                 updated_at timestamptz NOT NULL,
                                 CONSTRAINT user_roles_pkey PRIMARY KEY (user_id, permission_id)
);
CREATE INDEX idx_user_roles_fk_permission_uid ON user_permissions USING btree (permission_id);

ALTER TABLE user_roles ADD CONSTRAINT user_roles_role_id_fkey FOREIGN KEY (role_id) REFERENCES roles(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE user_permissions ADD CONSTRAINT user_roles_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE;

-- +migrate Down
DROP TABLE user_roles;