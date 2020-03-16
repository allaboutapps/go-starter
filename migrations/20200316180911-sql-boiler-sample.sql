-- +migrate Up
CREATE TABLE pilots (
    id uuid NOT NULL DEFAULT uuid_generate_v4 (),
    name text NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz,
    CONSTRAINT pilot_pkey PRIMARY KEY (id)
);

CREATE TABLE jets (
    id uuid NOT NULL DEFAULT uuid_generate_v4 (),
    pilot_id uuid NOT NULL,
    age integer NOT NULL,
    name text NOT NULL,
    color text NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz,
    CONSTRAINT jet_pkey PRIMARY KEY (id),
    CONSTRAINT jet_pilots_fkey FOREIGN KEY (pilot_id) REFERENCES pilots (id)
);

CREATE TABLE languages (
    id uuid NOT NULL DEFAULT uuid_generate_v4 (),
    created_at timestamptz NOT NULL,
    updated_at timestamptz,
    "language" text NOT NULL,
    CONSTRAINT language_pkey PRIMARY KEY (id)
);

CREATE TABLE pilot_languages (
    pilot_id uuid NOT NULL,
    language_id uuid NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz,
    CONSTRAINT pilot_language_pkey PRIMARY KEY (pilot_id, language_id),
    CONSTRAINT pilot_language_pilots_fkey FOREIGN KEY (pilot_id) REFERENCES pilots (id),
    CONSTRAINT pilot_language_languages_fkey FOREIGN KEY (language_id) REFERENCES languages (id)
);

-- +migrate Down
DROP TABLE pilot_languages;

DROP TABLE languages;

DROP TABLE jets;

DROP TABLE pilots;

