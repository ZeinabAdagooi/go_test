-- name: CreateUser
-- CreateUser creates a new user.
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20) UNIQUE NOT NULL,
    otp VARCHAR(6) NOT NULL,
    otp_expiration_time TIMESTAMP
);

