package db

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"storex/models"
)

func ListUsers(filters *models.UserFilterParams) ([]models.ListUsersResponse, error) {
	query := `
    SELECT
        u.id,
        u.name,
        u.email,
        u.phone,
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
		query += fmt.Sprintf(" AND (u.name ILIKE $%d OR u.email ILIKE $%d OR u.phone ILIKE $%d)", argIndex, argIndex+1, argIndex+2)
		args = append(args, "%"+filters.Search+"%", "%"+filters.Search+"%", "%"+filters.Search+"%")
		argIndex += 3
	}

	if len(filters.UserTypes) > 0 {
		query += fmt.Sprintf(" AND u.user_type = ANY($%d)", argIndex)
		args = append(args, pq.Array(filters.UserTypes))
		argIndex++
	}
	//if filters.UserType != "" {
	//	query += fmt.Sprintf(" AND u.user_type = $%d", argIndex)
	//	args = append(args, filters.UserType)
	//	argIndex++
	//}

	if len(filters.Roles) > 0 {
		query += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM user_roles ur2 WHERE ur2.user_id = u.id AND ur2.role = ANY($%d))", argIndex)
		args = append(args, pq.Array(filters.Roles))
		argIndex++
	}
	//if filters.Role != "" {
	//	query += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM user_roles ur2 WHERE ur2.user_id = u.id AND ur2.role = $%d)", argIndex)
	//	args = append(args, filters.Role)
	//	argIndex++
	//}

	if len(filters.AssetStatus) > 0 {
		query += fmt.Sprintf(`
			AND EXISTS (
				SELECT 1 FROM asset_status s2
				WHERE s2.assigned_to_user = u.id
				AND s2.status = ANY($%d)
				AND s2.archived_at IS NULL
			)`, argIndex)
		args = append(args, pq.Array(filters.AssetStatus))
		argIndex++
	}
	//if filters.AssetStatus != "" {
	//	query += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM asset_status s2 WHERE s2.assigned_to_user = u.id AND s2.status = $%d AND s2.archived_at IS NULL)", argIndex)
	//	args = append(args, filters.AssetStatus)
	//	argIndex++
	//}

	query += fmt.Sprintf(`
	GROUP BY u.id
	ORDER BY u.name
	LIMIT $%d OFFSET $%d
	`, argIndex, argIndex+1)

	args = append(args, filters.Limit, filters.Offset)
	//query += `
	//GROUP BY u.id
	//ORDER BY u.name
	//LIMIT $` + fmt.Sprint(argIndex) + ` OFFSET $` + fmt.Sprint(argIndex+1)
	//args = append(args, filters.Limit, filters.Offset)

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.ListUsersResponse
	for rows.Next() {
		var u models.ListUsersResponse
		var roles []sql.NullString
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.UserType, pq.Array(&roles), &u.AssignedAssetCount)
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

func SoftDeleteUser(tx *sql.Tx, userID string) error {
	_, err := tx.Exec(`
		UPDATE users SET archived_at = NOW() WHERE id = $1 AND archived_at IS NULL
	`, userID)

	if err != nil {
		return err
	}
	return nil
}

func GetUserDetailsByUserID(userID string) (models.UserDetails, error) {
	query := `
		SELECT
			u.id, u.name, u.email, u.user_type,
			COALESCE(array_agg(DISTINCT ur.role) FILTER (WHERE ur.role IS NOT NULL), '{}') AS roles,
			COUNT(DISTINCT ast.id) FILTER (
				WHERE ast.status = 'assigned'
				AND ast.assigned_to_user = u.id
				AND ast.archived_at IS NULL
			) AS assigned_asset_count
		FROM users u
		LEFT JOIN user_roles ur ON ur.user_id = u.id
		LEFT JOIN asset_status ast ON ast.assigned_to_user = u.id
		WHERE u.id = $1
		GROUP BY u.id
	`

	var user models.UserDetails
	var roles []sql.NullString
	err := DB.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.Email, &user.UserType, pq.Array(&roles), &user.AssignedAssetCount)
	if err != nil {
		return user, err
	}

	for _, r := range roles {
		if r.Valid {
			user.Roles = append(user.Roles, r.String)
		}
	}
	return user, nil
}
