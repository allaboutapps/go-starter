-- +migrate Up
CREATE TABLE app_user_profiles (
    user_id uuid NOT NULL,
    legal_accepted_at timestamptz,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CONSTRAINT app_user_profiles_pkey PRIMARY KEY (user_id)
);

ALTER TABLE app_user_profiles
    ADD CONSTRAINT app_user_profiles_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE;

-- +migrate Down
DROP TABLE IF EXISTS app_user_profiles;

