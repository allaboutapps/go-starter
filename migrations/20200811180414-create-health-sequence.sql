-- +migrate Up
CREATE SEQUENCE seq_health;

-- +migrate Down
DROP SEQUENCE IF EXISTS seq_health;

