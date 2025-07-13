CREATE TYPE owner_type AS ENUM('remote_state', 'client');

CREATE TABLE IF NOT EXISTS assets  (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    model_id UUID REFERENCES asset_models(id) NOT NULL,
    specs_id UUID NOT NULL,
    serial_no TEXT UNIQUE NOT NULL,
    owned_by owner_type NOT NULL,
    purchased_date TIMESTAMPTZ NOT NULL,
    warranty_start_date TIMESTAMPTZ,
    warranty_exp_date TIMESTAMPTZ,
    created_by UUID REFERENCES users(id) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    updated_by UUID REFERENCES users(id),
    archived_at TIMESTAMPTZ,
    archived_by UUID REFERENCES users(id)
);