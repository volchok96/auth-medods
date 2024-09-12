CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    user_guid UUID NOT NULL UNIQUE,
    ip VARCHAR(39),
    hashed_refresh_token TEXT,
    email VARCHAR(100) NOT NULL
);