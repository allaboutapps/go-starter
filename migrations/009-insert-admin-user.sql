-- +migrate Up
INSERT INTO permissions (id, scope, created_at, updated_at)
    VALUES ('ff150c28-ea0b-45ba-93a1-fc04ff78c2cd', 'cms', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
INSERT INTO permissions (id, scope, created_at, updated_at)
    VALUES ('ff3a2e29-8d7a-447a-93cb-ff1cfb7947a6', 'root', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
-- noinspection SqlIdentifierLength

INSERT INTO users (id, username, password, salt, created_at, updated_at) -- Z36ADo7Nba plain pw
    VALUES ('248f1e83-15b6-407b-933f-6afced1ef95b',
            'admin',
            '4b4f9a825579bcfae7107469d191a5f54e4ae47b060419a5c86fca0c6b0eb3206add2949080ee78772c83047e8fd875cbf3cdf997cae8e1089c7021460e527e3e014943bb86517444bcd314caf7a11f6f46e93e383f281f3ecd94d5346065538af090735c974478533ed796a39eed79ba6bd658fbd287609789f1008b1092730143ca70590e06dbd98182a5d6f03b01e3c25edf7680acc5f6e92bcd4a1056eae7b64df5a0045799851f83e6052af38487f408d656fae22c42d157efe946960213561dbbcfb24c98a226d314e2509156caa8154d312bd50b8211a4a335359d84f55ac7f8cffa7591c6f3f1aaf3d7e8a71cc8695eac0347ccd89d428ed7ce17702c882b4adbb83038dd7df6217441294daf3776edaed7ace624f6718232131473ae989760aee21c70dbb59c1e653be780a9f967051b550b2bdc2e4af4a606fc1a1bfc324c55aecc161a75d6e90ecc72870b26f8d6d2658ae3bb6429c873e6e26c970c487ed28f4695405d007684ee6dff26ee5fd3a257dd7f272e011213341a3455952f54c11a9b9f1551b70bdb933be0f140d4e2f8822c86f96765bc7b67a5e785c442e28b50a053b708ebcb84fb65f45546d84023bcbcaa4f330cae63908342064325223faef8f8e2f95aeccec254d1eccfea5167d240d6e56ea3ef1dfacfd9c36877a3116e67388601f030d5a7cb20a87cb8235355b08a471e32351f2ae4ad4',
            'd914c2d4fdc27deb940c4fcdfabba9f2f7ccb9cf58d598788cdce0f3e531b547ec9d64a35d56403932dabbc431dc3400c17d016b566d44c496f46c51665c46e4',
            CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
            );
INSERT INTO user_permissions (user_id, permission_id, created_at, updated_at)
    VALUES ('248f1e83-15b6-407b-933f-6afced1ef95b', 'ff150c28-ea0b-45ba-93a1-fc04ff78c2cd', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
INSERT INTO user_permissions (user_id, permission_id, created_at, updated_at)
    VALUES ('248f1e83-15b6-407b-933f-6afced1ef95b', 'ff3a2e29-8d7a-447a-93cb-ff1cfb7947a6', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- +migrate Down
DELETE FROM user_permissions WHERE user_id = '248f1e83-15b6-407b-933f-6afced1ef95b' AND permission_id = 'ff150c28-ea0b-45ba-93a1-fc04ff78c2cd';
DELETE FROM user_permissions WHERE user_id = '248f1e83-15b6-407b-933f-6afced1ef95b' AND permission_id = 'ff3a2e29-8d7a-447a-93cb-ff1cfb7947a6';
DELETE FROM users WHERE id = '248f1e83-15b6-407b-933f-6afced1ef95b';
DELETE FROM permissions WHERE id = 'ff150c28-ea0b-45ba-93a1-fc04ff78c2cd';
DELETE FROM permissions WHERE id = 'ff3a2e29-8d7a-447a-93cb-ff1cfb7947a6';