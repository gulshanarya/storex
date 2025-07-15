package handlers

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"storex/db"
	"storex/middleware"
	"storex/models"
	"strconv"
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

	err = db.InsertAssetStatus(tx, assetID, "available")
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

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		http.Error(w, "limit is not a number", http.StatusBadRequest)
		return
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		http.Error(w, "page is not a number", http.StatusBadRequest)
		return
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Parse query params
	params := models.ListAssetsQueryParams{
		Search:    r.URL.Query().Get("search"),
		AssetType: r.URL.Query().Get("asset_type"),
		Status:    r.URL.Query().Get("status"),
		OwnedBy:   r.URL.Query().Get("owned_by"),
		Limit:     limit,
		Offset:    (page - 1) * limit,
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

func UpdateAsset(w http.ResponseWriter, r *http.Request) {
	assetID := chi.URLParam(r, "id")

	var req models.UpdateAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// get existing asset (with model + asset_type)
	existingAsset, err := db.GetAssetWithModel(assetID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "asset not found", http.StatusNotFound)
		return
	}

	// Start transaction
	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "failed to begin transaction", http.StatusInternalServerError)
		return
	}
	defer db.TxFinalizer(tx, &err)

	userID := middleware.GetUserID(r)
	err = db.UpdateAsset(tx, &req, assetID, userID)

	if err != nil {
		http.Error(w, "failed to update asset", http.StatusInternalServerError)
		return
	}

	// Update specs if provided
	if req.Specs != nil {
		err := db.UpdateSpecsByType(tx, existingAsset.AssetType, existingAsset.SpecsID, req.Specs)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "failed to update asset(specs)", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Asset updated successfully",
		"asset_id": assetID,
	})
}

func AssignAsset(w http.ResponseWriter, r *http.Request) {
	var req models.AssignAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "could not begin transaction", http.StatusInternalServerError)
		return
	}
	defer db.TxFinalizer(tx, &err)

	// Check if asset exists and is available
	isAvailable, err := db.IsAssetAvailable(tx, req.AssetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !isAvailable {
		http.Error(w, "Asset is not available for assignment", http.StatusBadRequest)
		return
	}

	// Create asset_status entry
	if err := db.InsertAssetStatusToUser(tx, &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Asset assigned successfully"))
}

func RetrieveAsset(w http.ResponseWriter, r *http.Request) {
	assetID := chi.URLParam(r, "asset_id")
	if assetID == "" {
		http.Error(w, "asset ID required", http.StatusBadRequest)
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "could not begin transaction", http.StatusInternalServerError)
		return
	}
	defer db.TxFinalizer(tx, &err)

	// Check if asset is currently assigned
	statusID, err := db.GetActiveAssignedStatusID(tx, assetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set archived_at (retrieval time)
	if err := db.ArchiveAssetStatus(tx, statusID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Asset retrieved successfully"))
}

func AssetTimeline(w http.ResponseWriter, r *http.Request) {
	assetID := r.URL.Query().Get("asset_id")
	if assetID == "" {
		http.Error(w, "asset ID is required", http.StatusBadRequest)
		return
	}

	timelines, err := db.FetchAssetTimeline(assetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(timelines)
}

func UserAssetTimeline(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user ID is required", http.StatusBadRequest)
		return
	}

	timelines, err := db.FetchUserAssetTimeline(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(timelines)
}

func GetAssetsByUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")
	if userID == "" {
		http.Error(w, "user id required", http.StatusBadRequest)
	}

}
