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
		//self create user
		tx, err := db.DB.Begin()
		if err != nil {
			http.Error(w, "failed to start transaction", http.StatusInternalServerError)
			return
		}
		defer db.TxFinalizer(tx, &err)

		name := utils.ExtractNameFromEmail(user.Email)
		fullName := strings.Split(name, ".")

		for _, name := range fullName {
			user.Name += name
		}

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
