-- +migrate Up
CREATE TABLE applicant_tags (
    applicant_id uuid NOT NULL,
    tag_id uuid NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CONSTRAINT applicant_tags_pkey PRIMARY KEY (applicant_id, tag_id)
);

ALTER TABLE applicant_tags
    ADD CONSTRAINT applicant_tags_applicant_id_fkey FOREIGN KEY (applicant_id) REFERENCES applicants (id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE applicant_tags
    ADD CONSTRAINT applicant_tags_tag_id_fkey FOREIGN KEY (tag_id) REFERENCES tags (id) ON UPDATE CASCADE ON DELETE CASCADE;

-- +migrate Down
DROP TABLE IF EXISTS applicant_tags;

