CREATE TYPE user_type AS ENUM ('full_time', 'intern', 'freelancer');

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    phone TEXT UNIQUE,
    user_type user_type NOT NULL DEFAULT 'full_time',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES users(id),
    updated_at TIMESTAMPTZ,
    updated_by UUID REFERENCES users(id),
    archived_at TIMESTAMPTZ,
    archived_by UUID REFERENCES users(id)
);

-- indexes for text search
CREATE UNIQUE INDEX IF NOT EXISTS uniq_active_emails ON users(TRIM(LOWER(email))) WHERE archived_at IS NULL;
CREATE INDEX idx_users_lower_name ON users(LOWER(name));
CREATE INDEX idx_user_phone ON users(phone)

-- Index for filtering user_type
CREATE INDEX idx_users_user_type ON users(user_type);

