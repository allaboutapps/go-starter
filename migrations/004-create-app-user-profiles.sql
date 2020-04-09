-- +migrate Up
CREATE TABLE app_user_profiles (
                                 user_id uuid NOT NULL,
                                 has_gdpr_opt_out bool NOT NULL DEFAULT false,
                                 legal_accepted_at timestamptz NULL,
                                 created_at timestamptz NOT NULL,
                                 updated_at timestamptz NOT NULL,
                                 CONSTRAINT App_user_profiles_pkey PRIMARY KEY (user_id)
);
ALTER TABLE app_user_profiles ADD CONSTRAINT app_user_profiles_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE;

-- +migrate Down
DROP TABLE app_user_profiles;