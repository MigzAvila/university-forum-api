--  GIVE all users read permissions
INSERT INTO users_permissions
SELECT id, (SELECT id FROM permissions WHERE code = 'forums:read') FROM USERS;


-- GIVE SPECIFIC USER A permissions
INSERT INTO users_permissions(user_id, permission_id)
VALUES (
(SELECT id FROM users WHERE email = 'migz@example.com'),
(SELECT id FROM permissions WHERE code = 'forums:write')
);

-- List the activated users and their permisisons
SELECT email, array_agg(permissions.code) AS permisisons
FROM permissions
INNER JOIN users_permissions
ON users_permissions.permission_id = permissions.id
INNER JOIN users
ON users_permissions.user_id = users.id
WHERE users.activated = true
GROUP BY email;
