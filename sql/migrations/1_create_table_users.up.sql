CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(65) NOT NULL,
    first_name VARCHAR(128) NOT NULL,
    last_name VARCHAR(128) NULL,
    email VARCHAR(512) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL
);










