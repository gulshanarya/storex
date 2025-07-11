package db

import (
	"database/sql"
	"github.com/lib/pq"
	"storex/models"
)

func ListUsers() ([]models.UserDetails, error) {
	query := `
    SELECT
        u.id,
        u.name,
        u.email,
        u.user_type,
        COALESCE(array_agg(ur.role) FILTER (WHERE ur.role IS NOT NULL), '{}') AS roles,
        COUNT(DISTINCT s.id) FILTER (
            WHERE s.status = 'assigned'
            AND s.assigned_to_user = u.id
            AND s.archived_at IS NULL
        ) AS assigned_asset_count
    FROM users u
    LEFT JOIN user_roles ur ON ur.user_id = u.id
    LEFT JOIN asset_status s ON s.assigned_to_user = u.id
        AND s.status = 'assigned' AND s.archived_at IS NULL
    WHERE u.archived_at IS NULL
    GROUP BY u.id
    ORDER BY u.name;
    `

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserDetails
	for rows.Next() {
		var u models.UserDetails
		var roles []sql.NullString
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.UserType, pq.Array(&roles), &u.AssignedAssetCount)
		if err != nil {
			return nil, err
		}

		for _, r := range roles {
			if r.Valid {
				u.Roles = append(u.Roles, r.String)
			}
		}

		users = append(users, u)
	}
	return users, nil
}

func CreateProtectedUser(tx *sql.Tx, user *models.User) (string, error) {
	var userID string

	err := tx.QueryRow("INSERT INTO users(name, email, phone, user_type) VALUES ($1, $2, $3, $4) RETURNING id", user.Name, user.Email, user.Phone, user.UserType).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}
