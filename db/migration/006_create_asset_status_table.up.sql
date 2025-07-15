CREATE TYPE asset_status_type as ENUM('assigned', 'available', 'waiting_repair', 'service', 'damaged', 'disposed');
CREATE TABLE IF NOT EXISTS asset_status (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id UUID REFERENCES assets(id),
    status asset_status_type NOT NULL,
    assigned_to_user UUID REFERENCES users(id), --NULLABLE FIELD
    sent_to_service UUID REFERENCES services(id), --NULLABLE FIELD
    created_at TIMESTAMPTZ DEFAULT NOW(), //assigned time
    archived_at TIMESTAMPTZ //retrieved time
);

CREATE UNIQUE INDEX uniq_current_asset_status ON asset_status(asset_id)
    WHERE archived_at IS NULL;

CREATE INDEX idx_asset_status_assetid_status_active
    ON asset_status(asset_id, status)
    WHERE archived_at IS NULL;


CREATE INDEX idx_asset_status_user_status_active
    ON asset_status(assigned_to_user, status)
    WHERE archived_at IS NULL;
