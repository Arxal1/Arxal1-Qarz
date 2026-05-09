CREATE TABLE installments (
    id SERIAL PRIMARY KEY,
    event_id INT NOT NULL,
    due_date DATE NOT NULL,
    amount BIGINT NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    reminder_sent BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE
);