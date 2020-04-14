-- +migrate Up
CREATE TABLE role_application_states (
    role_id uuid NOT NULL,
    application_state_id uuid NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CONSTRAINT role_application_states_pkey PRIMARY KEY (role_id, application_state_id)
);

ALTER TABLE role_application_states
    ADD CONSTRAINT role_application_states_role_id_fkey FOREIGN KEY (role_id) REFERENCES roles (id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE role_application_states
    ADD CONSTRAINT role_application_states_application_state_id_fkey FOREIGN KEY (application_state_id) REFERENCES application_states (id) ON UPDATE CASCADE ON DELETE CASCADE;

-- +migrate Down
DROP TABLE application_states;

