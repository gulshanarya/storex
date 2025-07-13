package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"storex/models"
)

func GetOrCreateModel(tx *sql.Tx, brandID string, req *models.CreateModelRequest) (string, error) {
	var modelID string

	log.Println(req.Name, brandID, req.AssetType)
	querySelect := `
		SELECT id FROM asset_models
		WHERE LOWER(name) = LOWER($1) AND brand_id = $2 AND asset_type = $3
	`
	err := tx.QueryRow(querySelect, req.Name, brandID, req.AssetType).Scan(&modelID)
	if err == nil {
		return modelID, nil // Model exists
	}
	if err != sql.ErrNoRows {
		return "", err
	}

	// Insert new model
	queryInsert := `
		INSERT INTO asset_models (name, brand_id, asset_type)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err = tx.QueryRow(queryInsert, req.Name, brandID, req.AssetType).Scan(&modelID)
	if err != nil {
		return "", fmt.Errorf("failed to create model: %w", err)
	}

	return modelID, nil
}

func GetOrCreateBrand(tx *sql.Tx, brandName string) (string, error) {
	var brandID string

	querySelect := `SELECT id FROM asset_brands WHERE LOWER(name) = LOWER($1)`
	err := tx.QueryRow(querySelect, brandName).Scan(&brandID)
	if err == nil {
		return brandID, nil // Brand exists
	}

	if err != sql.ErrNoRows {
		return "", err
	}

	// Insert new brand
	queryInsert := `INSERT INTO asset_brands(name) VALUES ($1) RETURNING id`
	err = tx.QueryRow(queryInsert, brandName).Scan(&brandID)
	if err != nil {
		return "", fmt.Errorf("failed to create brand: %w", err)
	}

	return brandID, nil
}

func InsertSpecsAndReturnID(tx *sql.Tx, assetType string, specs interface{}) (string, error) {
	var query string
	var args []any
	var specsID string

	switch assetType {
	case "laptop":
		// Re-marshal into JSON
		specsBytes, err := json.Marshal(specs)
		if err != nil {
			return "", fmt.Errorf("failed to re-marshal specs: %w", err)
		}

		var s models.LaptopSpecs
		if err := json.Unmarshal(specsBytes, &s); err != nil {
			return "", fmt.Errorf("failed to unmarshal laptop specs: %w", err)
		}
		query = `INSERT INTO laptop_specs (processor, ram_gb, storage_gb, storage_type, screen_size_inch, has_charger) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
		args = []any{s.Processor, s.RAMGB, s.StorageGB, s.StorageType, s.ScreenSizeInch, s.HasCharger}

	case "mouse":
		// Re-marshal into JSON
		specsBytes, err := json.Marshal(specs)
		if err != nil {
			return "", fmt.Errorf("failed to re-marshal specs: %w", err)
		}

		var s models.MouseSpecs
		if err := json.Unmarshal(specsBytes, &s); err != nil {
			return "", fmt.Errorf("failed to unmarshal mouse specs: %w", err)
		}
		query = `INSERT INTO mouse_specs (type, dpi, number_of_buttons) VALUES ($1, $2, $3) RETURNING id`
		args = []any{s.Type, s.DPI, s.NumberOfButtons}

	case "monitor":
		// Re-marshal into JSON
		specsBytes, err := json.Marshal(specs)
		if err != nil {
			return "", fmt.Errorf("failed to re-marshal specs: %w", err)
		}

		var s models.MonitorSpecs
		if err := json.Unmarshal(specsBytes, &s); err != nil {
			return "", fmt.Errorf("failed to unmarshal monitor specs: %w", err)
		}
		query = `INSERT INTO monitor_specs (screen_size_inch, resolution, refresh_rate, panel_type) VALUES ($1, $2, $3, $4) RETURNING id`
		args = []any{s.ScreenSizeInch, s.Resolution, s.RefreshRate, s.PanelType}

	case "mobile":
		// Re-marshal into JSON
		specsBytes, err := json.Marshal(specs)
		if err != nil {
			return "", fmt.Errorf("failed to re-marshal specs: %w", err)
		}

		var s models.MobileSpecs
		if err := json.Unmarshal(specsBytes, &s); err != nil {
			return "", fmt.Errorf("failed to unmarshal mobile specs: %w", err)
		}
		query = `INSERT INTO mobile_specs (os, ram_gb, storage_gb, has_dual_sim) VALUES ($1, $2, $3, $4) RETURNING id`
		args = []any{s.OS, s.RAMGB, s.StorageGB, s.HasDualSIM}

	case "sim":
		// Re-marshal into JSON
		specsBytes, err := json.Marshal(specs)
		if err != nil {
			return "", fmt.Errorf("failed to re-marshal specs: %w", err)
		}

		var s models.SIMSpecs
		if err := json.Unmarshal(specsBytes, &s); err != nil {
			return "", fmt.Errorf("failed to unmarshal sim specs: %w", err)
		}
		query = `INSERT INTO sim_specs (carrier, phone_number, data_limit_gb) VALUES ($1, $2, $3) RETURNING id`
		args = []any{s.Carrier, s.PhoneNumber, s.DataLimitGB}

	case "hard_disk":
		// Re-marshal into JSON
		specsBytes, err := json.Marshal(specs)
		if err != nil {
			return "", fmt.Errorf("failed to re-marshal specs: %w", err)
		}

		var s models.HardDiskSpecs
		if err := json.Unmarshal(specsBytes, &s); err != nil {
			return "", fmt.Errorf("failed to unmarshal hard-disk specs: %w", err)
		}
		query = `INSERT INTO hard_disk_specs (capacity_gb, type) VALUES ($1, $2) RETURNING id`
		args = []any{s.CapacityGB, s.Type}

	case "pen_drive":
		// Re-marshal into JSON
		specsBytes, err := json.Marshal(specs)
		if err != nil {
			return "", fmt.Errorf("failed to re-marshal specs: %w", err)
		}

		var s models.PenDriveSpecs
		if err := json.Unmarshal(specsBytes, &s); err != nil {
			return "", fmt.Errorf("failed to unmarshal pen-drive specs: %w", err)
		}
		query = `INSERT INTO pen_drive_specs (capacity_gb, usb_version) VALUES ($1, $2) RETURNING id`
		args = []any{s.CapacityGB, s.USBVersion}

	case "accessories":
		// Re-marshal into JSON
		specsBytes, err := json.Marshal(specs)
		if err != nil {
			return "", fmt.Errorf("failed to re-marshal specs: %w", err)
		}

		var s models.AccessorySpecs
		if err := json.Unmarshal(specsBytes, &s); err != nil {
			return "", fmt.Errorf("failed to unmarshal accessory specs: %w", err)
		}
		query = `INSERT INTO accessory_specs (name, description, compatible_with) VALUES ($1, $2, $3) RETURNING id`
		args = []any{s.Name, s.Description, s.CompatibleWith}

	default:
		return "", fmt.Errorf("unsupported asset type: %s", assetType)
	}

	err := tx.QueryRow(query, args...).Scan(&specsID)
	if err != nil {
		return "", fmt.Errorf("failed to insert %s specs: %w", assetType, err)
	}

	return specsID, nil
}

func InsertAssetAndReturnID(tx *sql.Tx, modelID string, specsID string, req *models.CreateAssetRequest, authUserID string) (string, error) {
	query := `
		INSERT INTO assets (
			model_id, specs_id, serial_no, owned_by,
			purchased_date, warranty_start_date, warranty_exp_date,
			created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	var assetID string
	err := tx.QueryRow(query,
		modelID,
		specsID,
		req.SerialNo,
		req.OwnedBy,
		req.PurchasedDate,
		req.WarrantyStartDate,
		req.WarrantyExpDate,
		authUserID,
	).Scan(&assetID)

	if err != nil {
		return "", err
	}
	return assetID, nil
}

func InsertAssetStatus(tx *sql.Tx, assetID string, status string) error {

	statusInsertQuery := `
    INSERT INTO asset_status (asset_id, status)
    VALUES ($1, $2)
`
	_, err := tx.Exec(statusInsertQuery, assetID, status)
	if err != nil {
		return fmt.Errorf("failed to insert asset status: %w", err)
	}
	return nil
}

func ListAssets(params *models.ListAssetsQueryParams) ([]models.ListAssetsResponse, error) {
	// Base query
	query := `
			SELECT 
				a.id, a.serial_no, a.owned_by, a.purchased_date, 
				m.name AS model_name, m.asset_type, 
				b.name AS brand_name,
				s.status
			FROM assets a
			JOIN asset_models m ON a.model_id = m.id
			JOIN asset_brands b ON m.brand_id = b.id
			LEFT JOIN asset_status s ON s.asset_id = a.id AND s.archived_at IS NULL
			WHERE 1=1
		`

	var args []any
	argIndex := 1

	// Text search on brand/model/serial_no
	if params.Search != "" {
		query += fmt.Sprintf(` AND (
				b.name ILIKE $%d OR 
				m.name ILIKE $%d OR 
				a.serial_no ILIKE $%d)`, argIndex, argIndex+1, argIndex+2)
		args = append(args, "%"+params.Search+"%", "%"+params.Search+"%", "%"+params.Search+"%")
		argIndex += 3
	}

	// Filter asset_type
	if params.AssetType != "" {
		query += fmt.Sprintf(" AND m.asset_type = $%d", argIndex)
		args = append(args, params.AssetType)
		argIndex++
	}

	// Filter status
	if params.Status != "" {
		query += fmt.Sprintf(" AND s.status = $%d", argIndex)
		args = append(args, params.Status)
		argIndex++
	}

	// Filter owned_by
	if params.OwnedBy != "" {
		query += fmt.Sprintf(" AND a.owned_by = $%d", argIndex)
		args = append(args, params.OwnedBy)
		argIndex++
	}

	query += " ORDER BY a.created_at DESC"

	// Execute query
	rows, err := DB.Query(query, args...)
	if err != nil {
		log.Printf("ListAssets query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Parse results
	var assets []models.ListAssetsResponse
	for rows.Next() {
		var item models.ListAssetsResponse
		err := rows.Scan(&item.ID, &item.SerialNo, &item.OwnedBy, &item.PurchasedDate, &item.ModelName, &item.AssetType, &item.BrandName, &item.Status)
		if err != nil {
			log.Printf("Row scan error: %v", err)
			continue
		}
		assets = append(assets, item)
	}

	return assets, nil
}
