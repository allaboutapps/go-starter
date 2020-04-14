-- +migrate Up
CREATE TABLE applicants (
    id uuid NOT NULL,
    application_state_id uuid NOT NULL,
    name text NOT NULL,
    email text NOT NULL,
    job_category text,
    previous_projects text,
    start_date date NOT NULL,
    work_hours real NOT NULL,
    salary real NOT NULL,
    device text NOT NULL,
    discovery text NOT NULL,
    resume_url text NOT NULL,
    optional_files text[],
    gdpr_accepted boolean NOT NULL,
    last_change date NOT NULL,
    seniority text,
    first_interview date,
    second_interview date,
    notes text,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CONSTRAINT applicants_pkey PRIMARY KEY (id),
    CONSTRAINT applicants_email_key UNIQUE (email)
);

ALTER TABLE applicants
    ADD CONSTRAINT applicants_application_state_id_fkey FOREIGN KEY (application_state_id) REFERENCES application_states (id) ON UPDATE CASCADE ON DELETE CASCADE;

-- +migrate Down
DROP TABLE applicants;

