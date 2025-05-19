-- Create extension for UUID support
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Table: wallets
CREATE TABLE IF NOT EXISTS wallets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),       -- Unique wallet ID
    user_id UUID NOT NULL UNIQUE,                                -- ID of the user who owns the wallet
    balance BIGINT NOT NULL DEFAULT 0,                    -- Balance in smallest currency unit (e.g., cents)
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Table: transactions
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),       -- Unique transaction ID
    from_wallet UUID REFERENCES wallets(id) ON DELETE SET NULL,  -- Sender's wallet (nullable for deposits)
    to_wallet UUID REFERENCES wallets(id) ON DELETE SET NULL,    -- Receiver's wallet (nullable for withdrawals)
    amount BIGINT NOT NULL CHECK (amount > 0),            -- Transaction amount must be positive
    type VARCHAR(20) NOT NULL CHECK (type IN ('deposit', 'withdrawal', 'transfer')),  -- Transaction type
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_wallet_user_id ON wallets(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_wallets ON transactions(from_wallet, to_wallet);
