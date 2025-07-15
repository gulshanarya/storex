package models

type User struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Phone    *string `json:"phone"`
	Role     string  `json:"role"`
	UserType string  `json:"user_type"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"` // optional: you can return a new refresh token or reuse
}

type ListUsersResponse struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Email              string   `json:"email"`
	Phone              *string  `json:"phone"`
	Roles              []string `json:"roles"`
	UserType           string   `json:"user_type"`
	AssignedAssetCount int      `json:"assigned_asset_count"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	UserType *string `json:"user_type"`
}

type CreateBrandRequest struct {
	Name string `json:"name"`
}

type CreateModelRequest struct {
	Name      string             `json:"name"`
	AssetType string             `json:"asset_type"`
	Brand     CreateBrandRequest `json:"brand"`
}
