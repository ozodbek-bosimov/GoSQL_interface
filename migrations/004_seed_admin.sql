-- Default admin user: login=admin, password=admin123
-- Password hashed with bcrypt cost 10
INSERT INTO users (name, login, password, role)
VALUES (
    'Administrator',
    'admin',
    '$2a$10$.xi0p6IzBOAUMHM7XUj9U.X6pPgWgeG73ILhOQSxbFK2qXxbRsRx6',
    'admin'
)
ON CONFLICT (login) DO UPDATE SET password = EXCLUDED.password;
