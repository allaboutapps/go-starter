-- +migrate Up
CREATE TABLE role_notification_templates (
                                 role_id uuid NOT NULL,
                                 notification_template_id uuid NOT NULL,
                                 created_at timestamptz NOT NULL,
                                 updated_at timestamptz NOT NULL,
                                 CONSTRAINT role_notification_templates_pkey PRIMARY KEY (role_id, notification_template_id)
);

ALTER TABLE role_notification_templates ADD CONSTRAINT role_notification_templates_role_id_fkey FOREIGN KEY (role_id) REFERENCES roles(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE role_notification_templates ADD CONSTRAINT role_notification_templates_notification_template_id_fkey FOREIGN KEY (notification_template_id) REFERENCES notificaton_templates(id) ON UPDATE CASCADE ON DELETE CASCADE;

-- +migrate Down
DROP TABLE role_notification_templates;