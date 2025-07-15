package db

import (
	"database/sql"
	"storex/models"
)

func TxFinalizer(tx *sql.Tx, err *error) {
	if p := recover(); p != nil {
		tx.Rollback()
		panic(p)
	} else if *err != nil {
		tx.Rollback()
	} else {
		*err = tx.Commit()
	}
}

func UpdateLaptopSpecsByID(tx *sql.Tx, specsID string, specs *models.LaptopSpecsUpdate) error {
	return buildSpecsUpdateSQL(tx, "laptop_specs", specsID, map[string]interface{}{
		"processor":        specs.Processor,
		"ram_gb":           specs.RAMGB,
		"storage_gb":       specs.StorageGB,
		"storage_type":     specs.StorageType,
		"screen_size_inch": specs.ScreenSizeInch,
		"has_charger":      specs.HasCharger,
	})
}

func UpdateMouseSpecsByID(tx *sql.Tx, specsID string, specs *models.MouseSpecsUpdate) error {
	return buildSpecsUpdateSQL(tx, "mouse_specs", specsID, map[string]interface{}{
		"type":              specs.Type,
		"dpi":               specs.DPI,
		"number_of_buttons": specs.NumberOfButtons,
	})
}

func UpdateMonitorSpecsByID(tx *sql.Tx, specsID string, specs *models.MonitorSpecsUpdate) error {
	return buildSpecsUpdateSQL(tx, "monitor_specs", specsID, map[string]interface{}{
		"screen_size_inch": specs.ScreenSizeInch,
		"resolution":       specs.Resolution,
		"refresh_rate":     specs.RefreshRate,
		"panel_type":       specs.PanelType,
	})
}

func UpdateMobileSpecsByID(tx *sql.Tx, specsID string, specs *models.MobileSpecsUpdate) error {
	return buildSpecsUpdateSQL(tx, "mobile_specs", specsID, map[string]interface{}{
		"os":           specs.OS,
		"ram_gb":       specs.RAMGB,
		"storage_gb":   specs.StorageGB,
		"has_dual_sim": specs.HasDualSIM,
	})
}

func UpdateSIMSpecsByID(tx *sql.Tx, specsID string, specs *models.SIMSpecsUpdate) error {
	return buildSpecsUpdateSQL(tx, "sim_specs", specsID, map[string]interface{}{
		"carrier":       specs.Carrier,
		"phone_number":  specs.PhoneNumber,
		"data_limit_gb": specs.DataLimitGB,
	})
}

func UpdateHardDiskSpecsByID(tx *sql.Tx, specsID string, specs *models.HardDiskSpecsUpdate) error {
	return buildSpecsUpdateSQL(tx, "hard_disk_specs", specsID, map[string]interface{}{
		"capacity_gb": specs.CapacityGB,
		"type":        specs.Type,
	})
}

func UpdatePenDriveSpecsByID(tx *sql.Tx, specsID string, specs *models.PenDriveSpecsUpdate) error {
	return buildSpecsUpdateSQL(tx, "pen_drive_specs", specsID, map[string]interface{}{
		"capacity_gb": specs.CapacityGB,
		"usb_version": specs.USBVersion,
	})
}

func UpdateAccessorySpecsByID(tx *sql.Tx, specsID string, specs *models.AccessorySpecsUpdate) error {
	return buildSpecsUpdateSQL(tx, "accessory_specs", specsID, map[string]interface{}{
		"name":            specs.Name,
		"description":     specs.Description,
		"compatible_with": specs.CompatibleWith,
	})
}
