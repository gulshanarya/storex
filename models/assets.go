package models

import (
	"time"
)

type ListAssetsResponse struct {
	ID            string `json:"id"`
	SerialNo      string `json:"serial_no"`
	OwnedBy       string `json:"owned_by"`
	PurchasedDate string `json:"purchased_date"`
	ModelName     string `json:"model_name"`
	AssetType     string `json:"asset_type"`
	BrandName     string `json:"brand_name"`
	Status        string `json:"status"`
}

type ListAssetsQueryParams struct {
	Search    string
	AssetType string
	Status    string
	OwnedBy   string
	Limit     int
	Offset    int
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

type UpdateAssetRequest struct {
	Model             *CreateModelRequest `json:"model"`
	SerialNo          *string             `json:"serial_no"`
	OwnedBy           *string             `json:"owned_by"` // ENUM: "remote_state", "client"
	PurchasedDate     *time.Time          `json:"purchased_date"`
	WarrantyStartDate *time.Time          `json:"warranty_start_date"`
	WarrantyExpDate   *time.Time          `json:"warranty_exp_date"`
	Specs             interface{}         `json:"specs"` // Raw specs for dynamic routing
	Status            *string             `json:"status"`
}

type AssetWithModel struct {
	ID                string
	ModelID           string
	SpecsID           string
	SerialNo          string
	OwnedBy           string
	PurchasedDate     time.Time
	WarrantyStartDate *time.Time
	WarrantyExpDate   *time.Time
	AssetType         string
}

type AssignAssetRequest struct {
	AssetID string `json:"asset_id"`
	UserID  string `json:"user_id"`
}

type AssetTimeline struct {
	Status         string     `json:"status"`
	AssignedToUser *string    `json:"assigned_to_user,omitempty"`
	SentToService  *string    `json:"sent_to_service,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	ArchivedAt     *time.Time `json:"archived_at,omitempty"`
}
