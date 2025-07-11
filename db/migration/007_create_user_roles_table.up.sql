CREATE TYPE user_role AS ENUM ('admin', 'asset_manager', 'employee_manager', 'employee');

CREATE TABLE user_roles (
    user_id UUID REFERENCES users(id),
    role user_role NOT NULL,
    PRIMARY KEY(user_id, role)
)