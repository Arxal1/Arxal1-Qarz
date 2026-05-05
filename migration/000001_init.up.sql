-- migrations/000001_init.up.sql

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    phone VARCHAR(20) UNIQUE,
    full_name VARCHAR(255),
    role VARCHAR(20) DEFAULT 'client', -- 'owner', 'manager', 'client'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_identified BOOLEAN DEFAULT FALSE,
    pinfl VARCHAR(14) UNIQUE
);