CREATE TABLE laptop_specs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    processor TEXT NOT NULL,
    ram_gb INTEGER NOT NULL,
    storage_gb INTEGER NOT NULL,
    storage_type TEXT NOT NULL, -- 'HDD' or 'SSD'
    screen_size_inch NUMERIC NOT NULL,
    has_charger BOOLEAN DEFAULT true
);

CREATE TABLE mouse_specs (
     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
     type TEXT NOT NULL, -- 'wired', 'wireless'
     dpi INTEGER,
     number_of_buttons INTEGER
);

CREATE TABLE monitor_specs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    screen_size_inch NUMERIC NOT NULL,
    resolution TEXT, -- e.g. "1920x1080"
    refresh_rate INTEGER,
    panel_type TEXT -- 'IPS', 'VA', 'TN'
);

CREATE TABLE hard_disk_specs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    capacity_gb INTEGER NOT NULL,
    type TEXT NOT NULL -- 'HDD', 'SSD'
);

CREATE TABLE pen_drive_specs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    capacity_gb INTEGER NOT NULL,
    usb_version TEXT -- e.g. "2.0", "3.0"
);

CREATE TABLE mobile_specs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    os TEXT NOT NULL,      -- e.g. "Android", "iOS"
    ram_gb INTEGER,
    storage_gb INTEGER,
    has_dual_sim BOOLEAN
);

CREATE TABLE sim_specs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    carrier TEXT NOT NULL,     -- e.g. "Jio", "Airtel"
    phone_number TEXT UNIQUE,
    data_limit_gb INTEGER
);

CREATE TABLE accessory_specs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    compatible_with TEXT -- free-text info, e.g. "laptop, mobile"
);







