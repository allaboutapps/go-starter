-- +migrate Up
CREATE TABLE roles (
                              id uuid NOT NULL,
                              role text NULL,
                              created_at timestamptz NOT NULL,
                              updated_at timestamptz NOT NULL,
                              CONSTRAINT roles_pkey PRIMARY KEY (id)
);

-- +migrate Down
DROP TABLE roles;