-- +migrate Up
CREATE TABLE users (
    id uuid NOT NULL DEFAULT uuid_generate_v4 (),
    username varchar(255),
    password text,
    google_id varchar(255) DEFAULT NULL::character varying,
    google_info text,
    facebook_id varchar(255) DEFAULT NULL::character varying,
    facebook_info text,
    is_active bool NOT NULL DEFAULT TRUE,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CONSTRAINT users_facebook_id_key UNIQUE (facebook_id),
    CONSTRAINT users_google_id_key UNIQUE (google_id),
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT users_username_key UNIQUE (username)
);

-- +migrate Down
DROP TABLE users;

