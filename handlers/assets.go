package handlers

import (
	"log"
	"net/http"
	"storex/db"
	"storex/middleware"
	"storex/models"
)

func CreateAsset(w http.ResponseWriter, r *http.Request) {
	var req models.CreateAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.SerialNo == "" || req.PurchasedDate == nil {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "failed to begin transaction", http.StatusInternalServerError)
		return
	}
	defer db.TxFinalizer(tx, &err)

	// Get or create brand
	brandID, err := db.GetOrCreateBrand(tx, req.Model.Brand.Name)
	if err != nil {
		http.Error(w, "failed to resolve brand", http.StatusInternalServerError)
		return
	}

	// Get or create model under the brand
	modelID, err := db.GetOrCreateModel(tx, brandID, &req.Model)
	if err != nil {
		http.Error(w, "failed to resolve model", http.StatusInternalServerError)
		return
	}

	// Insert into specs table based on asset_type
	specsID, err := db.InsertSpecsAndReturnID(tx, req.Model.AssetType, req.Specs)
	if err != nil {
		http.Error(w, "failed to insert specs: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert into assets table
	authUserID := middleware.GetUserID(r)
	log.Println(authUserID)

	//insert into asset table
	assetID, err := db.InsertAssetAndReturnID(tx, modelID, specsID, &req, authUserID)

	if err != nil {
		http.Error(w, "failed to insert asset: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.InsertAssetStatus(tx, assetID, req.Status)
	if err != nil {
		http.Error(w, "failed to insert asset status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Asset created successfully",
		"asset_id": assetID,
	})
}

func ListAssets(w http.ResponseWriter, r *http.Request) {
	// Parse query params
	params := models.ListAssetsQueryParams{
		Search:    r.URL.Query().Get("search"),
		AssetType: r.URL.Query().Get("asset_type"),
		Status:    r.URL.Query().Get("status"),
		OwnedBy:   r.URL.Query().Get("owned_by"),
	}

	assets, err := db.ListAssets(&params)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "failed to list assets", http.StatusInternalServerError)
		return
	}

	if len(assets) == 0 {
		http.Error(w, "no assets found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(assets)
}
