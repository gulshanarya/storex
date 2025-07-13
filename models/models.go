package models

import (
	"database/sql"
	"time"
)

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

type UserDetails struct {
	ID                 string         `json:"id"`
	Name               string         `json:"name"`
	Email              string         `json:"email"`
	Phone              sql.NullString `json:"phone"`
	Roles              []string       `json:"roles"`
	UserType           string         `json:"user_type"`
	AssignedAssetCount int            `json:"assigned_asset_count"`
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
type CreateAssetRequest struct {
	Model             CreateModelRequest `json:"model"`
	SerialNo          string             `json:"serial_no"`
	OwnedBy           string             `json:"owned_by"` // ENUM: "remote_state", "client"
	PurchasedDate     *time.Time         `json:"purchased_date"`
	WarrantyStartDate *time.Time         `json:"warranty_start_date"`
	WarrantyExpDate   *time.Time         `json:"warranty_exp_date"`
	Specs             interface{}        `json:"specs"` // Raw specs for dynamic routing
	Status            string             `json:"status"`
}
