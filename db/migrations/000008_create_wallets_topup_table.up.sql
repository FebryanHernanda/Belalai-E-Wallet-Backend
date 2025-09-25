CREATE TABLE wallets_topup (
    wallets_id INT NOT NULL,
    topup_id INT NOT NULL,
    PRIMARY KEY (wallets_id, topup_id),
    CONSTRAINT fk_wallets_topup_wallets FOREIGN KEY (wallets_id) REFERENCES wallets(id),
    CONSTRAINT fk_wallets_topup_topup FOREIGN KEY (topup_id) REFERENCES topup(id)
);