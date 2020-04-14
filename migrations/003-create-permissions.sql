-- +migrate Up
CREATE TABLE permissions (
    id uuid NOT NULL,
    scope varchar(255) NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CONSTRAINT permissions_pkey PRIMARY KEY (id),
    CONSTRAINT permissions_scope_key UNIQUE (scope)
);

-- +migrate Down
DROP TABLE permissions;

