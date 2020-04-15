-- Don't forget to maintain order here, foreign keys!
DROP TABLE IF EXISTS video_tags;

DROP TABLE IF EXISTS tags;

DROP TABLE IF EXISTS videos;

DROP TABLE IF EXISTS sponsors;

DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS type_monsters;

DROP TYPE IF EXISTS workday;

CREATE TYPE workday AS enum (
    'monday',
    'tuesday',
    'wednesday',
    'thursday',
    'friday'
);

DROP DOMAIN IF EXISTS uint3;

CREATE DOMAIN uint3 AS numeric CHECK (value >= 0
    AND value < power(2::numeric, 3::numeric));

CREATE TABLE users (
    id serial PRIMARY KEY NOT NULL
);

CREATE TABLE sponsors (
    id serial PRIMARY KEY NOT NULL
);

CREATE TABLE videos (
    id serial PRIMARY KEY NOT NULL,
    user_id int NOT NULL,
    sponsor_id int UNIQUE,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (sponsor_id) REFERENCES sponsors (id)
);

CREATE TABLE tags (
    id serial PRIMARY KEY NOT NULL
);

CREATE TABLE video_tags (
    video_id int NOT NULL,
    tag_id int NOT NULL,
    PRIMARY KEY (video_id, tag_id),
    FOREIGN KEY (video_id) REFERENCES videos (id),
    FOREIGN KEY (tag_id) REFERENCES tags (id)
);

DROP TYPE IF EXISTS my_int_array;

CREATE DOMAIN my_int_array AS int[];

CREATE TABLE type_monsters (
    id serial PRIMARY KEY NOT NULL,
    enum_use workday NOT NULL,
    bool_zero bool,
    bool_one bool NULL,
    bool_two bool NOT NULL,
    bool_three bool NULL DEFAULT FALSE,
    bool_four bool NULL DEFAULT TRUE,
    bool_five bool NOT NULL DEFAULT FALSE,
    bool_six bool NOT NULL DEFAULT TRUE,
    string_zero varchar(1),
    string_one varchar(1) NULL,
    string_two varchar(1) NOT NULL,
    string_three varchar(1) NULL DEFAULT 'a',
    string_four varchar(1) NOT NULL DEFAULT 'b',
    string_five varchar(1000),
    string_six varchar(1000) NULL,
    string_seven varchar(1000) NOT NULL,
    string_eight varchar(1000) NULL DEFAULT 'abcdefgh',
    string_nine varchar(1000) NOT NULL DEFAULT 'abcdefgh',
    string_ten varchar(1000) NULL DEFAULT '',
    string_eleven varchar(1000) NOT NULL DEFAULT '',
    nonbyte_zero char(1),
    nonbyte_one char(1) NULL,
    nonbyte_two char(1) NOT NULL,
    nonbyte_three char(1) NULL DEFAULT 'a',
    nonbyte_four char(1) NOT NULL DEFAULT 'b',
    nonbyte_five char(1000),
    nonbyte_six char(1000) NULL,
    nonbyte_seven char(1000) NOT NULL,
    nonbyte_eight char(1000) NULL DEFAULT 'a',
    nonbyte_nine char(1000) NOT NULL DEFAULT 'b',
    byte_zero "char",
    byte_one "char" NULL,
    byte_two "char" NULL DEFAULT 'a',
    byte_three "char" NOT NULL,
    byte_four "char" NOT NULL DEFAULT 'b',
    big_int_zero bigint,
    big_int_one bigint NULL,
    big_int_two bigint NOT NULL,
    big_int_three bigint NULL DEFAULT 111111,
    big_int_four bigint NOT NULL DEFAULT 222222,
    big_int_five bigint NULL DEFAULT 0,
    big_int_six bigint NOT NULL DEFAULT 0,
    int_zero int,
    int_one int NULL,
    int_two int NOT NULL,
    int_three int NULL DEFAULT 333333,
    int_four int NOT NULL DEFAULT 444444,
    int_five int NULL DEFAULT 0,
    int_six int NOT NULL DEFAULT 0,
    float_zero decimal,
    float_one numeric,
    float_two numeric(2, 1),
    float_three numeric(2, 1),
    float_four numeric(2, 1) NULL,
    float_five numeric(2, 1) NOT NULL,
    float_six numeric(2, 1) NULL DEFAULT 1.1,
    float_seven numeric(2, 1) NOT NULL DEFAULT 1.1,
    float_eight numeric(2, 1) NULL DEFAULT 0.0,
    float_nine numeric(2, 1) NULL DEFAULT 0.0,
    bytea_zero bytea,
    bytea_one bytea NULL,
    bytea_two bytea NOT NULL,
    bytea_three bytea NOT NULL DEFAULT 'a',
    bytea_four bytea NULL DEFAULT 'b',
    bytea_five bytea NOT NULL DEFAULT 'abcdefghabcdefghabcdefgh',
    bytea_six bytea NULL DEFAULT 'hgfedcbahgfedcbahgfedcba',
    bytea_seven bytea NOT NULL DEFAULT '',
    bytea_eight bytea NOT NULL DEFAULT '',
    time_zero timestamp,
    time_one date,
    time_two timestamp NULL DEFAULT NULL,
    time_three timestamp NULL,
    time_four timestamp NOT NULL,
    time_five timestamp NULL DEFAULT '1999-01-08 04:05:06.789',
    time_six timestamp NULL DEFAULT '1999-01-08 04:05:06.789 -8:00',
    time_seven timestamp NULL DEFAULT 'January 8 04:05:06 1999 PST',
    time_eight timestamp NOT NULL DEFAULT '1999-01-08 04:05:06.789',
    time_nine timestamp NOT NULL DEFAULT '1999-01-08 04:05:06.789 -8:00',
    time_ten timestamp NOT NULL DEFAULT 'January 8 04:05:06 1999 PST',
    time_eleven date NULL,
    time_twelve date NOT NULL,
    time_thirteen date NULL DEFAULT '1999-01-08',
    time_fourteen date NULL DEFAULT 'January 8, 1999',
    time_fifteen date NULL DEFAULT '19990108',
    time_sixteen date NOT NULL DEFAULT '1999-01-08',
    time_seventeen date NOT NULL DEFAULT 'January 8, 1999',
    time_eighteen date NOT NULL DEFAULT '19990108',
    uuid_zero uuid,
    uuid_one uuid NULL,
    uuid_two uuid NULL DEFAULT NULL,
    uuid_three uuid NOT NULL,
    uuid_four uuid NULL DEFAULT '6ba7b810-9dad-11d1-80b4-00c04fd430c8',
    uuid_five uuid NOT NULL DEFAULT '6ba7b810-9dad-11d1-80b4-00c04fd430c8',
    integer_default integer DEFAULT '5' ::integer,
    varchar_default varchar(1000) DEFAULT 5::varchar,
    timestamp_notz timestamp without time zone DEFAULT (now() at time zone 'utc'),
    timestamp_tz timestamp with time zone DEFAULT (now() at time zone 'utc'),
    interval_nnull interval NOT NULL DEFAULT '21 days',
    interval_null interval NULL DEFAULT '23 hours',
    json_null json NULL,
    json_nnull json NOT NULL,
    jsonb_null jsonb NULL,
    jsonb_nnull jsonb NOT NULL,
    box_null box NULL,
    box_nnull box NOT NULL,
    cidr_null cidr NULL,
    cidr_nnull cidr NOT NULL,
    circle_null circle NULL,
    circle_nnull circle NOT NULL,
    double_prec_null double precision NULL,
    double_prec_nnull double precision NOT NULL,
    inet_null inet NULL,
    inet_nnull inet NOT NULL,
    line_null line NULL,
    line_nnull line NOT NULL,
    lseg_null lseg NULL,
    lseg_nnull lseg NOT NULL,
    macaddr_null macaddr NULL,
    macaddr_nnull macaddr NOT NULL,
    money_null money NULL,
    money_nnull money NOT NULL,
    path_null path NULL,
    path_nnull path NOT NULL,
    pg_lsn_null pg_lsn NULL,
    pg_lsn_nnull pg_lsn NOT NULL,
    point_null point NULL,
    point_nnull point NOT NULL,
    polygon_null polygon NULL,
    polygon_nnull polygon NOT NULL,
    tsquery_null tsquery NULL,
    tsquery_nnull tsquery NOT NULL,
    tsvector_null tsvector NULL,
    tsvector_nnull tsvector NOT NULL,
    txid_null txid_snapshot NULL,
    txid_nnull txid_snapshot NOT NULL,
    xml_null xml NULL,
    xml_nnull xml NOT NULL,
    intarr_null integer[] NULL,
    intarr_nnull integer[] NOT NULL,
    boolarr_null boolean[] NULL,
    boolarr_nnull boolean[] NOT NULL,
    varchararr_null varchar[] NULL,
    varchararr_nnull varchar[] NOT NULL,
    decimalarr_null decimal[] NULL,
    decimalarr_nnull decimal[] NOT NULL,
    byteaarr_null bytea[] NULL,
    byteaarr_nnull bytea[] NOT NULL,
    jsonbarr_null jsonb[] NULL,
    jsonbarr_nnull jsonb[] NOT NULL,
    jsonarr_null json[] NULL,
    jsonarr_nnull json[] NOT NULL,
    customarr_null my_int_array NULL,
    customarr_nnull my_int_array NOT NULL,
    domainuint3_nnull uint3 NOT NULL
);

