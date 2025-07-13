package models

type ListAssetsQueryParams struct {
	Search    string
	AssetType string
	Status    string
	OwnedBy   string
}

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
