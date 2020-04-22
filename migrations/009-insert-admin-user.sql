-- +migrate Up
INSERT INTO permissions (id, scope, created_at, updated_at)
    VALUES ('ff150c28-ea0b-45ba-93a1-fc04ff78c2cd', 'cms', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO permissions (id, scope, created_at, updated_at)
    VALUES ('ff3a2e29-8d7a-447a-93cb-ff1cfb7947a6', 'root', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- noinspection SqlIdentifierLength
INSERT INTO users (id, username, PASSWORD, created_at, updated_at) -- Z36ADo7Nba plain pw
    VALUES ('248f1e83-15b6-407b-933f-6afced1ef95b', 'admin', '$argon2id$v=19$m=65536,t=1,p=4$YpwpZ1S9tYczmqcq9iRCVA$SD/ft4GbHFAkmgulclfqZa/yq9FeTlpOxBGEUVUVPzI', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO user_permissions (user_id, permission_id, created_at, updated_at)
    VALUES ('248f1e83-15b6-407b-933f-6afced1ef95b', 'ff150c28-ea0b-45ba-93a1-fc04ff78c2cd', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO user_permissions (user_id, permission_id, created_at, updated_at)
    VALUES ('248f1e83-15b6-407b-933f-6afced1ef95b', 'ff3a2e29-8d7a-447a-93cb-ff1cfb7947a6', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- +migrate Down
DELETE FROM user_permissions
WHERE user_id = '248f1e83-15b6-407b-933f-6afced1ef95b'
    AND permission_id = 'ff150c28-ea0b-45ba-93a1-fc04ff78c2cd';

DELETE FROM user_permissions
WHERE user_id = '248f1e83-15b6-407b-933f-6afced1ef95b'
    AND permission_id = 'ff3a2e29-8d7a-447a-93cb-ff1cfb7947a6';

DELETE FROM users
WHERE id = '248f1e83-15b6-407b-933f-6afced1ef95b';

DELETE FROM permissions
WHERE id = 'ff150c28-ea0b-45ba-93a1-fc04ff78c2cd';

DELETE FROM permissions
WHERE id = 'ff3a2e29-8d7a-447a-93cb-ff1cfb7947a6';

