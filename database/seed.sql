-- Seed data for development/testing
-- Run AFTER schema.sql

INSERT INTO users (email, first_name, last_name, otp_code)
VALUES
    ('john@example.com', 'John', 'Doe', '123456'),
    ('jane@example.com', 'Jane', 'Smith', '654321')
ON CONFLICT (email) DO NOTHING;
