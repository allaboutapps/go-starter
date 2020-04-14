-- +migrate Up
CREATE TABLE users (
    id uuid NOT NULL,
    username varchar(255) NULL,
    password varchar(2048) NULL,
    salt varchar(2048) NULL,
    google_id varchar(255) NULL DEFAULT NULL::character varying,
    google_info text NULL,
    facebook_id varchar(255) NULL DEFAULT NULL::character varying,
    facebook_info text NULL,
    is_active bool NULL DEFAULT TRUE,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CONSTRAINT users_facebook_id_key UNIQUE (facebook_id),
    CONSTRAINT users_google_id_key UNIQUE (google_id),
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT users_username_key UNIQUE (username)
);

-- +migrate Down
DROP TABLE users;

