-- +migrate Up
CREATE TABLE notification_templates (
    id uuid NOT NULL,
    text timestamptz NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CONSTRAINT notification_templates_pkey PRIMARY KEY (id)
);

-- +migrate Down
DROP TABLE notification_templates;

