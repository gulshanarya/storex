CREATE TYPE asset_status_type as ENUM('assigned', 'available', 'waiting_repair', 'service', 'damaged', 'disposed');
CREATE TABLE IF NOT EXISTS asset_status (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id UUID REFERENCES assets(id),
    status asset_status_type NOT NULL,
    assigned_to_user UUID REFERENCES users(id), --NULLABLE FIELD
    sent_to_service UUID REFERENCES services(id), --NULLABLE FIELD
    created_at TIMESTAMPTZ DEFAULT NOW(),
    archived_at TIMESTAMPTZ
)