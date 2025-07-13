package models

type UserFilterParams struct {
	Search      string `json "search"`
	UserType    string `json "user_type"`
	Role        string `json "role"`
	AssetStatus string `json "asset_status"`
	Limit       int    `json "limit"`
	Offset      int    `json "offset"`
}
