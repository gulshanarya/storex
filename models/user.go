package models

import "time"

type UserFilterParams struct {
	Search      string   `json "search"`
	UserTypes   []string `json "user_type"`
	Roles       []string `json "role"`
	AssetStatus []string `json "asset_status"`
	Limit       int      `json "limit"`
	Offset      int      `json "offset"`
}

type AssignedAsset struct {
	AssetID    string    `json:"asset_id"`
	ModelName  string    `json:"model_name"`
	BrandName  string    `json:"brand_name"`
	Status     string    `json:"status"`
	AssignedAt time.Time `json:"assigned_at"`
}

type UserDetails struct {
	ID                 string          `json:"id"`
	Name               string          `json:"name"`
	Email              string          `json:"email"`
	UserType           string          `json:"user_type"`
	Roles              []string        `json:"roles"`
	AssignedAssetCount int             `json:"asset_status"`
	AssignedAssets     []AssignedAsset `json:"assigned_assets"`
}
