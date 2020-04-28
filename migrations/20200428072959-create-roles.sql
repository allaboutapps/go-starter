-- +migrate Up
CREATE TABLE roles (
    id uuid NOT NULL DEFAULT uuid_generate_v4 (),
    role text NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CONSTRAINT roles_pkey PRIMARY KEY (id)
);

-- +migrate Down
DROP TABLE IF EXISTS roles;

