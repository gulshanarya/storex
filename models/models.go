package models

import "database/sql"

type User struct {
	Id       string         `json:"id"`
	Name     string         `json:"name"`
	Email    string         `json:"email"`
	Phone    sql.NullString `json:"phone"`
	Role     string         `json:"role"`
	UserType string         `json:"user_type"`
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
