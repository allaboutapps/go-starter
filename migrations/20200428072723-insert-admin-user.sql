-- +migrate Up
INSERT INTO users (id, username, PASSWORD, scopes, created_at, updated_at) -- Z36ADo7Nba plain pw
    VALUES ('248f1e83-15b6-407b-933f-6afced1ef95b', 'admin', '$argon2id$v=19$m=65536,t=1,p=4$YpwpZ1S9tYczmqcq9iRCVA$SD/ft4GbHFAkmgulclfqZa/yq9FeTlpOxBGEUVUVPzI', '{cms}', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- +migrate Down
DELETE FROM users
WHERE id = '248f1e83-15b6-407b-933f-6afced1ef95b';

