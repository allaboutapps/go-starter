-- +migrate Up
CREATE TABLE tags (
                              id uuid NOT NULL,
                              tag text NULL,
                              created_at timestamptz NOT NULL,
                              updated_at timestamptz NOT NULL,
                              CONSTRAINT tags_pkey PRIMARY KEY (id)
);

-- +migrate Down
DROP TABLE tags;