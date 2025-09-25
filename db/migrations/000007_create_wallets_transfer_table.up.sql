CREATE TABLE wallets_transfer (
    wallets_id INT NOT NULL,
    transfer_id INT NOT NULL,
    PRIMARY KEY (wallets_id, transfer_id),
    CONSTRAINT fk_wallets_transfer_wallets FOREIGN KEY (wallets_id) REFERENCES wallets(id),
    CONSTRAINT fk_wallets_transfer_transfer FOREIGN KEY (transfer_id) REFERENCES transfer(id)
);