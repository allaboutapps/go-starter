-- +migrate Up
CREATE TYPE user_scope AS ENUM (
    'app',
    'cms'
);

-- +migrate Down
DROP TYPE IF EXISTS user_scope;

