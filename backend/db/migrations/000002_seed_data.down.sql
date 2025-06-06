-- Remove admin user
DELETE FROM users WHERE email = 'admin@example.com';

-- Remove all services
DELETE FROM services; 