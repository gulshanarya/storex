package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"storex/models"
	"strings"
	"time"
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

	query += ` ORDER BY a.created_at DESC
	LIMIT $` + fmt.Sprint(argIndex) + ` OFFSET $` + fmt.Sprint(argIndex+1)

	args = append(args, params.Limit, params.Offset)

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

func GetAssetWithModel(assetID string) (*models.AssetWithModel, error) {
	query := `
		SELECT
			a.id,
			a.model_id,
			a.specs_id,
			a.serial_no,
			a.owned_by,
			a.purchased_date,
			a.warranty_start_date,
			a.warranty_exp_date,
			m.asset_type
		FROM assets a
		JOIN asset_models m ON a.model_id = m.id
		WHERE a.id = $1 AND a.archived_at IS NULL
	`

	var asset models.AssetWithModel
	err := DB.QueryRow(query, assetID).Scan(
		&asset.ID,
		&asset.ModelID,
		&asset.SpecsID,
		&asset.SerialNo,
		&asset.OwnedBy,
		&asset.PurchasedDate,
		&asset.WarrantyStartDate,
		&asset.WarrantyExpDate,
		&asset.AssetType,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &asset, nil
}

func UpdateAsset(tx *sql.Tx, req *models.UpdateAssetRequest, assetID string, userID string) error {
	// Dynamically build asset update query
	setClauses := []string{}
	args := []interface{}{}
	argID := 1

	if req.SerialNo != nil {
		setClauses = append(setClauses, fmt.Sprintf("serial_no = $%d", argID))
		args = append(args, *req.SerialNo)
		argID++
	}

	if req.OwnedBy != nil {
		setClauses = append(setClauses, fmt.Sprintf("owned_by = $%d", argID))
		args = append(args, *req.OwnedBy)
		argID++
	}

	if req.PurchasedDate != nil {
		setClauses = append(setClauses, fmt.Sprintf("purchased_date = $%d", argID))
		args = append(args, *req.PurchasedDate)
		argID++
	}

	if req.WarrantyStartDate != nil {
		setClauses = append(setClauses, fmt.Sprintf("warranty_start_date = $%d", argID))
		args = append(args, *req.WarrantyStartDate)
		argID++
	}

	if req.WarrantyExpDate != nil {
		setClauses = append(setClauses, fmt.Sprintf("warranty_exp_date = $%d", argID))
		args = append(args, *req.WarrantyExpDate)
		argID++
	}

	// Add updated_by and updated_at
	setClauses = append(setClauses,
		fmt.Sprintf("updated_by = $%d", argID),
		fmt.Sprintf("updated_at = $%d", argID+1),
	)
	args = append(args, userID, time.Now())
	argID += 2

	// Finalize query
	if len(setClauses) > 0 {
		query := fmt.Sprintf(`
			UPDATE assets SET %s WHERE id = $%d
		`, strings.Join(setClauses, ", "), argID)
		args = append(args, assetID)

		if _, err := tx.Exec(query, args...); err != nil {
			return err
		}
	}

	return nil
}

func UpdateSpecsByType(tx *sql.Tx, assetType string, specsID string, specs interface{}) error {
	// Re-marshal into JSON
	specsBytes, err := json.Marshal(specs)
	if err != nil {
		return fmt.Errorf("failed to re-marshal specs: %w", err)
	}
	switch assetType {
	case "laptop":
		var s models.LaptopSpecsUpdate
		if err := json.Unmarshal(specsBytes, &s); err != nil {
			return fmt.Errorf("failed to unmarshal laptop specs: %w", err)
		}
		return UpdateLaptopSpecsByID(tx, specsID, &s)

	case "mouse":
		var s models.MouseSpecsUpdate
		if err := json.Unmarshal(specsBytes, &s); err != nil {
			return fmt.Errorf("failed to unmarshal specs: %w", err)
		}
		return UpdateMouseSpecsByID(tx, specsID, &s)

	case "monitor":
		var s models.MonitorSpecsUpdate
		if err := json.Unmarshal(specsBytes, &s); err != nil {
			return fmt.Errorf("failed to unmarshal specs: %w", err)
		}
		return UpdateMonitorSpecsByID(tx, specsID, &s)

	case "mobile":
		var s models.MobileSpecsUpdate
		if err := json.Unmarshal(specsBytes, &s); err != nil {
			return fmt.Errorf("failed to unmarshal specs: %w", err)
		}
		return UpdateMobileSpecsByID(tx, specsID, &s)

	case "sim":
		var s models.SIMSpecsUpdate
		if err := json.Unmarshal(specsBytes, &s); err != nil {
			return fmt.Errorf("failed to unmarshal specs: %w", err)
		}
		return UpdateSIMSpecsByID(tx, specsID, &s)

	case "hard_disk":
		var s models.HardDiskSpecsUpdate
		if err := json.Unmarshal(specsBytes, &s); err != nil {
			return fmt.Errorf("failed to unmarshal specs: %w", err)
		}
		return UpdateHardDiskSpecsByID(tx, specsID, &s)

	case "pen_drive":
		var s models.PenDriveSpecsUpdate
		if err := json.Unmarshal(specsBytes, &s); err != nil {
			return fmt.Errorf("failed to unmarshal specs: %w", err)
		}
		return UpdatePenDriveSpecsByID(tx, specsID, &s)

	case "accessories":
		var s models.AccessorySpecsUpdate
		if err := json.Unmarshal(specsBytes, &s); err != nil {
			return fmt.Errorf("failed to unmarshal specs: %w", err)
		}
		return UpdateAccessorySpecsByID(tx, specsID, &s)

	default:
		return errors.New("unsupported asset type: " + assetType)
	}
}

func buildSpecsUpdateSQL(tx *sql.Tx, table string, id string, fields map[string]interface{}) error {
	var sets []string
	var args []interface{}
	argPos := 1

	for col, val := range fields {
		if val != nil {
			sets = append(sets, fmt.Sprintf("%s = COALESCE($%d, %s)", col, argPos, col))
			args = append(args, val)
			argPos++
		}
	}

	if len(sets) == 0 {
		return nil // no updates
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", table, strings.Join(sets, ", "), argPos)
	log.Println(sets, args)
	_, err := tx.Exec(query, args...)
	return err
}

func IsAssetAvailable(tx *sql.Tx, assetID string) (bool, error) {
	query := `SELECT COUNT(*) FROM asset_status WHERE asset_id = $1 AND status = 'assigned' AND archived_at IS NULL`
	var count int
	err := tx.QueryRow(query, assetID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func InsertAssetStatusToUser(tx *sql.Tx, status *models.AssignAssetRequest) error {
	query := `INSERT INTO asset_status (asset_id, status, assigned_to_user) VALUES ($1, $2, $3)`
	_, err := tx.Exec(query, status.AssetID, "assigned", status.UserID)
	return err
}

func GetActiveAssignedStatusID(tx *sql.Tx, assetID string) (string, error) {
	query := `SELECT id FROM asset_status WHERE asset_id = $1 AND status = 'assigned' AND archived_at IS NULL`
	var id string
	err := tx.QueryRow(query, assetID).Scan(&id)
	if err == sql.ErrNoRows {
		return "", errors.New("No active assignment found for asset")
	} else if err != nil {
		return "", err
	}
	return id, nil
}

func ArchiveAssetStatus(tx *sql.Tx, statusID string) error {
	query := `UPDATE asset_status SET archived_at = NOW() WHERE id = $1`
	_, err := tx.Exec(query, statusID)
	return err
}

func FetchAssetTimeline(assetID string) ([]models.AssetTimeline, error) {
	query := `
		SELECT status, assigned_to_user, sent_to_service, created_at, archived_at
		FROM asset_status
		WHERE asset_id = $1
		ORDER BY created_at ASC
	`

	rows, err := DB.Query(query, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var timeline []models.AssetTimeline
	for rows.Next() {
		var t models.AssetTimeline
		var assignedTo sql.NullString
		var sentToService sql.NullString
		var archivedAt sql.NullTime
		if err := rows.Scan(&t.Status, &assignedTo, &sentToService, &t.CreatedAt, &archivedAt); err != nil {
			return nil, err
		}
		if assignedTo.Valid {
			t.AssignedToUser = &assignedTo.String
		}
		if sentToService.Valid {
			t.SentToService = &sentToService.String
		}
		if archivedAt.Valid {
			t.ArchivedAt = &archivedAt.Time
		}
		timeline = append(timeline, t)
	}
	return timeline, nil
}

func FetchUserAssetTimeline(userID string) ([]models.AssetTimeline, error) {
	//sent_to_service redundant remove later
	query := `
		SELECT status, asset_id, sent_to_service, created_at, archived_at
		FROM asset_status
		WHERE assigned_to_user = $1
		ORDER BY created_at ASC
	`

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var timeline []models.AssetTimeline
	for rows.Next() {
		var t models.AssetTimeline
		var assetID string
		var sentToService sql.NullString
		var archivedAt sql.NullTime
		if err := rows.Scan(&t.Status, &assetID, &sentToService, &t.CreatedAt, &archivedAt); err != nil {
			return nil, err
		}
		t.AssignedToUser = &userID
		if sentToService.Valid {
			t.SentToService = &sentToService.String
		}
		if archivedAt.Valid {
			t.ArchivedAt = &archivedAt.Time
		}
		timeline = append(timeline, t)
	}
	return timeline, nil
}

func NumberOfAssetsAssigned(tx *sql.Tx, userID string) (int, error) {
	assignedCount := 0
	err := tx.QueryRow(`
		SELECT COUNT(*) FROM asset_status
		WHERE assigned_to_user = $1 AND archived_at IS NULL
	`, userID).Scan(&assignedCount)

	if err != nil {
		return 0, err
	}
	return assignedCount, nil
}
