CREATE TABLE IF NOT EXISTS accounts
(
    id VARCHAR(32) NOT NULL PRIMARY KEY,
    account_id VARCHAR(32) NOT NULL,
    name VARCHAR(50) NOT NULL,
    birth_date TIMESTAMP WITH TIME ZONE NOT NULL,
    gender BOOLEAN NOT NULL,
    address TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT account_id_unique UNIQUE (account_id)
);

CREATE TABLE IF NOT EXISTS account_balances
(
    id VARCHAR(32) NOT NULL PRIMARY KEY,
    account_id VARCHAR(32) NOT NULL UNIQUE,
    amount DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT account_balance_unique UNIQUE (account_id),
    FOREIGN KEY (account_id) REFERENCES accounts(account_id)
);

CREATE TABLE IF NOT EXISTS account_transactions
(
    id VARCHAR(32) NOT NULL PRIMARY KEY,
    transaction_timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    amount DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    account_src_id VARCHAR(32) NOT NULL,
    account_dst_id VARCHAR(32) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (account_src_id) REFERENCES accounts(account_id),
    FOREIGN KEY (account_dst_id) REFERENCES accounts(account_id)
);