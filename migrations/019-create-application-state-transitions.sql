-- +migrate Up
CREATE TABLE application_state_transitions (
    applicant_id uuid NOT NULL,
    from_application_state_id uuid NOT NULL,
    to_application_state_id uuid NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CONSTRAINT application_state_transitions_pkey PRIMARY KEY (applicant_id, from_application_state_id, to_application_state_id)
);

ALTER TABLE application_state_transitions
    ADD CONSTRAINT application_state_transitions_applicant_id_fkey FOREIGN KEY (applicant_id) REFERENCES applicants (id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE application_state_transitions
    ADD CONSTRAINT application_state_transitions_from_application_state_id_fkey FOREIGN KEY (from_application_state_id) REFERENCES application_states (id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE application_state_transitions
    ADD CONSTRAINT application_state_transitions_to_application_state_id_fkey FOREIGN KEY (to_application_state_id) REFERENCES application_states (id) ON UPDATE CASCADE ON DELETE CASCADE;

-- +migrate Down
DROP TABLE application_state_transitions;

