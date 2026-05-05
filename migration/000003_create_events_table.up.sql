CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    client_id INT NOT NULL,
    event_type VARCHAR(20) NOT NULL, -- shipment (отгрузка в долг) или 'payment' (погашение долга)
    amount BIGINT NOT NULL,
    status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'confirmed', 'rejected'
    description TEXT,
    created_by INT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
);