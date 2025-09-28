CREATE TYPE topup_status AS ENUM ('success', 'failed', 'pending');
CREATE TABLE topup (
    id SERIAL PRIMARY KEY,
    amount INT NOT NULL,
    tax INT,
    payment_id INT NOT NULL,
    topup_status topup_status NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
    CONSTRAINT fk_topup_payment FOREIGN KEY (payment_id) REFERENCES payment_method(id)
);