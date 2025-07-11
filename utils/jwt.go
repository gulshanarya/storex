package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func ValidateJWT(tokenString string) (string, string, error) {
	jwtKey := []byte(os.Getenv("jwt_secret_key"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure token is signed with expected algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtKey, nil
	})

	if err != nil {
		return "", "", err
	}

	// Validate claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check expiration
		exp, ok := claims["exp"].(float64)
		if !ok || time.Now().Unix() > int64(exp) {
			return "", "", errors.New("token expired")
		}

		// Extract user_id
		userID, ok := claims["user_id"].(string)
		if !ok {
			return "", "", errors.New("user_id missing in token")
		}

		// Extract roles
		role, ok := claims["role"].(string)
		if !ok {
			return "", "", errors.New("role missing or invalid in token")
		}

		return userID, role, nil
	}

	return "", "", errors.New("invalid token")
}

func GenerateRefreshJWT(userID string) (string, error) {
	jwtKey := []byte(os.Getenv("jwt_secret_key"))
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().UTC().Add(7 * 24 * time.Hour),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func GenerateAccessJWT(userID string, role string) (string, error) {
	jwtKey := []byte(os.Getenv("jwt_secret_key"))
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().UTC().Add(15 * time.Minute).Unix(), // expires in 15mins
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
