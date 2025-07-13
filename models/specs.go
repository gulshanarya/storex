package models

// LaptopSpecs represents specs for a laptop asset
type LaptopSpecs struct {
	Processor      string  `json:"processor"`
	RAMGB          int     `json:"ram_gb"`
	StorageGB      int     `json:"storage_gb"`
	StorageType    string  `json:"storage_type"`
	ScreenSizeInch float64 `json:"screen_size_inch"`
	HasCharger     bool    `json:"has_charger"`
}

// MouseSpecs represents specs for a mouse asset
type MouseSpecs struct {
	Type            string `json:"type"`
	DPI             int    `json:"dpi"`
	NumberOfButtons int    `json:"number_of_buttons"`
}

// MonitorSpecs represents specs for a monitor asset
type MonitorSpecs struct {
	ScreenSizeInch float64 `json:"screen_size_inch"`
	Resolution     string  `json:"resolution"`
	RefreshRate    int     `json:"refresh_rate"`
	PanelType      string  `json:"panel_type"`
}

// MobileSpecs represents specs for a mobile asset
type MobileSpecs struct {
	OS         string `json:"os"`
	RAMGB      int    `json:"ram_gb"`
	StorageGB  int    `json:"storage_gb"`
	HasDualSIM bool   `json:"has_dual_sim"`
}

// SIMspecs represents specs for a SIM asset
type SIMSpecs struct {
	Carrier     string `json:"carrier"`
	PhoneNumber string `json:"phone_number"`
	DataLimitGB int    `json:"data_limit_gb"`
}

// HardDiskSpecs represents specs for a hard disk asset
type HardDiskSpecs struct {
	CapacityGB int    `json:"capacity_gb"`
	Type       string `json:"type"`
}

// PenDriveSpecs represents specs for a pen drive asset
type PenDriveSpecs struct {
	CapacityGB int    `json:"capacity_gb"`
	USBVersion string `json:"usb_version"`
}

// AccessorySpecs represents specs for generic accessories
type AccessorySpecs struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	CompatibleWith string `json:"compatible_with"`
}
