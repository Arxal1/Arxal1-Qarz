CREATE TABLE businesses (
        id SERIAL PRIMARY KEY,
        owner_id INT NOT NULL,
        name VARCHAR(255) NOT NULL,
        phone VARCHAR(20),
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

        FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE

);

CREATE TABLE clients (
    id SERIAL PRIMARY KEY,
    businnes_id INT NOT NULL,
    user_id INT,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    debt_limit BIGINT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (businnes_id) REFERENCES businesses(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,

    UNIQUE (businnes_id, phone)
);

