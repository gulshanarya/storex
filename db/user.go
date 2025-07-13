package db

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"storex/models"
)

func ListUsers(filters *models.UserFilterParams) ([]models.UserDetails, error) {
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
    LEFT JOIN asset_status s ON s.assigned_to_user = u.id AND s.archived_at IS NULL
    WHERE u.archived_at IS NULL
`

	var args []interface{}
	argIndex := 1

	if filters.Search != "" {
		query += fmt.Sprintf(" AND (LOWER(u.name) LIKE LOWER($%d) OR LOWER(u.email) LIKE LOWER($%d))", argIndex, argIndex+1)
		args = append(args, "%"+filters.Search+"%", "%"+filters.Search+"%")
		argIndex += 2
	}

	if filters.UserType != "" {
		query += fmt.Sprintf(" AND u.user_type = $%d", argIndex)
		args = append(args, filters.UserType)
		argIndex++
	}

	if filters.Role != "" {
		query += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM user_roles ur2 WHERE ur2.user_id = u.id AND ur2.role = $%d)", argIndex)
		args = append(args, filters.Role)
		argIndex++
	}

	if filters.AssetStatus != "" {
		query += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM asset_status s2 WHERE s2.assigned_to_user = u.id AND s2.status = $%d AND s2.archived_at IS NULL)", argIndex)
		args = append(args, filters.AssetStatus)
		argIndex++
	}

	query += `
    GROUP BY u.id
    ORDER BY u.name
    LIMIT $` + fmt.Sprint(argIndex) + ` OFFSET $` + fmt.Sprint(argIndex+1)
	args = append(args, filters.Limit, filters.Offset)

	rows, err := DB.Query(query, args...)
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

//func ListUsers() ([]models.UserDetails, error) {
//	query := `
//    SELECT
//        u.id,
//        u.name,
//        u.email,
//        u.user_type,
//        COALESCE(array_agg(ur.role) FILTER (WHERE ur.role IS NOT NULL), '{}') AS roles,
//        COUNT(DISTINCT s.id) FILTER (
//            WHERE s.status = 'assigned'
//            AND s.assigned_to_user = u.id
//            AND s.archived_at IS NULL
//        ) AS assigned_asset_count
//    FROM users u
//    LEFT JOIN user_roles ur ON ur.user_id = u.id
//    LEFT JOIN asset_status s ON s.assigned_to_user = u.id
//        AND s.status = 'assigned' AND s.archived_at IS NULL
//    WHERE u.archived_at IS NULL
//    GROUP BY u.id
//    ORDER BY u.name;
//    `
//
//	rows, err := DB.Query(query)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var users []models.UserDetails
//	for rows.Next() {
//		var u models.UserDetails
//		var roles []sql.NullString
//		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.UserType, pq.Array(&roles), &u.AssignedAssetCount)
//		if err != nil {
//			return nil, err
//		}
//
//		for _, r := range roles {
//			if r.Valid {
//				u.Roles = append(u.Roles, r.String)
//			}
//		}
//
//		users = append(users, u)
//	}
//	return users, nil
//}

func CreateProtectedUser(tx *sql.Tx, user *models.User) (string, error) {
	var userID string

	err := tx.QueryRow("INSERT INTO users(name, email, phone, user_type) VALUES ($1, $2, $3, $4) RETURNING id", user.Name, user.Email, user.Phone, user.UserType).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func UpdateUser(authUserID string, userID string, req *models.UpdateUserRequest) error {
	query := `
		UPDATE users SET
			name = COALESCE($1, name),
			email = COALESCE($2, email),
			phone = COALESCE($3, phone),
			user_type = COALESCE($4, user_type),
			updated_at = CURRENT_TIMESTAMP,
			updated_by = $5
		WHERE id = $6 AND archived_at IS NULL
	`

	_, err := DB.Exec(query,
		req.Name,
		req.Email,
		req.Phone,
		req.UserType,
		authUserID,
		userID,
	)
	if err != nil {
		return err
	}
	return nil
}
