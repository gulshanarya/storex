package db

import (
	"database/sql"
	"storex/models"
)

func CreateUser(tx *sql.Tx, user *models.User) (string, error) {
	var userID string

	err := tx.QueryRow("INSERT INTO users(name, email) VALUES ($1, $2) RETURNING id", user.Name, user.Email).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func CreateRole(tx *sql.Tx, userID string, role string) error {
	_, err := tx.Exec("INSERT INTO user_roles(user_id, role) VALUES ($1, $2)", userID, role)

	if err != nil {
		return err
	}
	return nil
}

func IsUserExist(email string) (bool, error) {

	err := DB.QueryRow(`
SELECT id FROM users WHERE email = $1 AND archived_at IS NULL`, email).Scan()

	if err != nil {
		return false, err
	}

	return true, nil
}

func GetUserDetails(user *models.User) error {
	// fetching detail of given user
	err := DB.QueryRow(`SELECT 
    u.id, u.name, u.email, u.phone, u.user_type,
    ur.role as role_name
	FROM users u
	LEFT JOIN user_roles ur ON ur.user_id = u.id
	WHERE u.email = $1 AND u.archived_at IS NULL;
`, user.Email).
		Scan(&user.Id, &user.Name, &user.Email, &user.Phone, &user.UserType, &user.Role)
	if err != nil {
		return err
	}
	return nil
}
