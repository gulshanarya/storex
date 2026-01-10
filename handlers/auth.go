package handlers

import (
	"database/sql"
	jsoniter "github.com/json-iterator/go"
	"log"
	"net/http"
	"storex/db"
	"storex/models"
	"storex/utils"
	"strings"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	if !utils.IsValidEmail(user.Email) {
		http.Error(w, "invalid email", http.StatusBadRequest)
		return
	}

	err = db.GetUserDetails(&user)
	if err == sql.ErrNoRows {
		if user.Role != "employee" {
			http.Error(w, "user with the given role not exists", http.StatusBadRequest)
			return
		}
		//self create user
		tx, err := db.DB.Begin()
		if err != nil {
			http.Error(w, "failed to start transaction", http.StatusInternalServerError)
			return
		}
		defer db.TxFinalizer(tx, &err)

		name := utils.ExtractNameFromEmail(user.Email)
		fullName := strings.Split(name, ".")

		//for _, name := range fullName {
		user.Name = fullName[0] + " " + fullName[1]
		//}

		userID, err := db.CreateUser(tx, &user)

		if err != nil {
			log.Println(err.Error())
			http.Error(w, "failed to create user", http.StatusInternalServerError)
			return
		}

		user.Role = "employee"
		err = db.CreateRole(tx, userID, user.Role)

		if err != nil {
			log.Println(err.Error())
			http.Error(w, "failed to create role", http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		log.Println("testing", err.Error())
		http.Error(w, "failed to find or create user", http.StatusInternalServerError)
		return
	}

	//	give login access
	accessToken, err := utils.GenerateAccessJWT(user.Id, user.Role)
	log.Println(err)
	if err != nil {
		http.Error(w, "failed to generate login access token", http.StatusInternalServerError)
		return
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshJWT(user.Id)
	if err != nil {
		http.Error(w, "failed to generate login refresh token", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	})
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
		AccessToken  string `json:"access_token" validate:"required"` // Expired access token (to get role)
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// validate expired jwt
	userID, role, err := utils.ValidateExpiredAccessJWT(req.AccessToken)
	if err != nil {
		log.Println("access token", err.Error())
		http.Error(w, "invalid access token", http.StatusUnauthorized)
		return
	}

	_, err = utils.ValidateRefreshJWT(req.RefreshToken)
	if err != nil {
		log.Println("refresh token", err.Error())
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	newAccessToken, err := utils.GenerateAccessJWT(userID, role)
	if err != nil {
		http.Error(w, "Failed to generate new access token", http.StatusInternalServerError)
		return
	}

	newRefreshToken, err := utils.GenerateRefreshJWT(userID)
	if err != nil {
		http.Error(w, "Failed to generate new refresh token", http.StatusInternalServerError)
		return
	}

	resp := models.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken, // no rotation
	}

	json.NewEncoder(w).Encode(resp)
}
