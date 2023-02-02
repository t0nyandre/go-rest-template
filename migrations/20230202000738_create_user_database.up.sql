CREATE TABLE users
(
    id VARCHAR(25) PRIMARY KEY,
    name TEXT,
    username VARCHAR(128) NOT NULL,
    password TEXT NOT NULL,
    email VARCHAR(254) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);