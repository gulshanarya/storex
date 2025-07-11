CREATE TYPE asset_type AS ENUM('laptop', 'mouse', 'monitor', 'hard_disk', 'pen_drive', 'mobile', 'sim', 'accessories');
CREATE TABLE IF NOT EXISTS asset_models (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    model_name TEXT NOT NULL,
    brand_id UUID REFERENCES asset_brands(id),
    asset_type asset_type NOT NULL
)