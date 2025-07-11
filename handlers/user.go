package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"storex/db"
	"storex/middleware"
	"storex/models"
	"storex/utils"
	"strings"
)

func ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := db.ListUsers()
	if err != nil {
		http.Error(w, "failed in getting users", http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(users)
}

func IsValidPhone(phone string) bool {
	if len(phone) == 10 {
		return true
	}
	return false
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

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

	if !utils.IsValidRole(user.Role) {
		http.Error(w, "invalid role", http.StatusBadRequest)
		return
	}

	if !utils.IsValidUserType(user.UserType) {
		http.Error(w, "invalid user type", http.StatusBadRequest)
		return
	}

	if user.Phone.Valid && !IsValidPhone(user.Phone.String) {
		http.Error(w, "invalid phone", http.StatusBadRequest)
		return
	}

	creatorRole := middleware.GetUserRole(r)

	if creatorRole == "admin" && user.Role == "admin" {
		http.Error(w, "Admins cannot create another admin", http.StatusForbidden)
		return
	}

	if creatorRole == "employee_manager" && user.Role != "employee" {
		http.Error(w, "Employee managers can only create employees", http.StatusForbidden)
		return
	}

	if creatorRole == "asset_manager" {
		http.Error(w, "Asset managers are not allowed to create users", http.StatusForbidden)
		return
	}

	isUserExist, err := db.IsUserExist(user.Email)
	if isUserExist {
		http.Error(w, "user with this email already exists", http.StatusInternalServerError)
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "failed in checking user", http.StatusInternalServerError)
		return
	}

	//create user
	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer db.TxFinalizer(tx, &err)

	name := utils.ExtractNameFromEmail(user.Email)
	fullName := strings.Split(name, ".")

	//skipping users given name instead take from mail
	user.Name = ""
	for _, name := range fullName {
		user.Name += name
	}

	userID, err := db.CreateProtectedUser(tx, &user)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	err = db.CreateRole(tx, userID, user.Role)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "failed to create role", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]string{
		"message": "User created successfully",
	})
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	if !utils.IsValidEmail(user.Email) {
	}
}
