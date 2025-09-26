CREATE TYPE transfer_status_enum AS ENUM ('success', 'failed', 'pending');
CREATE TABLE transfer (
    id SERIAL PRIMARY KEY,
    sender_wallet_id INT NOT NULL,
    receiver_wallet_id INT NOT NULL,
    amount INT NOT NULL,
    transfer_status transfer_status_enum NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_by_sender BOOLEAN DEFAULT FALSE,
    deleted_by_receiver BOOLEAN DEFAULT FALSE
);