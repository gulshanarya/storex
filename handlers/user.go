package handlers

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"storex/db"
	"storex/middleware"
	"storex/models"
	"storex/utils"
	"strconv"
	"strings"
)

func ListUsers(w http.ResponseWriter, r *http.Request) {

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		http.Error(w, "limit is not a number", http.StatusBadRequest)
		return
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		http.Error(w, "page is not a number", http.StatusBadRequest)
		return
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	parseMulti := func(param string) []string {
		values := strings.Split(r.URL.Query().Get(param), ",")
		var cleaned []string
		for _, v := range values {
			if trimmed := strings.TrimSpace(v); trimmed != "" {
				cleaned = append(cleaned, trimmed)
			}
		}
		return cleaned
	}

	// Parse query params
	params := models.UserFilterParams{
		Search:      r.URL.Query().Get("search"),
		UserTypes:   parseMulti("user_type"),
		AssetStatus: parseMulti("status"),
		Roles:       parseMulti("role"),
		Limit:       limit,
		Offset:      (page - 1) * limit,
	}

	//later will add multiple filter of a type
	users, err := db.ListUsers(&params)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "failed to list users", http.StatusInternalServerError)
		return
	}

	if len(users) == 0 {
		http.Error(w, "no users found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(users)
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

	//setting empty string as nil
	if user.Phone != nil && strings.TrimSpace(*user.Phone) == "" {
		user.Phone = nil
	}

	if user.Phone != nil && !utils.IsValidPhone(*user.Phone) {
		http.Error(w, "invalid phone", http.StatusBadRequest)
		return
	}

	creatorRole := middleware.GetUserRole(r)

	if creatorRole == "admin" && user.Role == "admin" {
		http.Error(w, "Admins cannot create another admin", http.StatusForbidden)
		return
	}

	if creatorRole == "employee_manager" && user.Role != "employee" {
		http.Error(w, "employee managers can only create employees", http.StatusForbidden)
		return
	}

	if creatorRole == "asset_manager" {
		http.Error(w, "Asset managers are not allowed to create users", http.StatusForbidden)
		return
	}

	isUserExist, err := db.IsUserExist(user.Email)
	if err != nil {
		http.Error(w, "failed in checking user existence", http.StatusInternalServerError)
		return
	}

	if isUserExist {
		http.Error(w, "user with this email already exists", http.StatusInternalServerError)
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
	userID := chi.URLParam(r, "user_id")
	if userID == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate provided fields
	if req.Email != nil && !utils.IsValidEmail(*req.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}
	if req.UserType != nil && !utils.IsValidUserType(*req.UserType) {
		http.Error(w, "Invalid user_type", http.StatusBadRequest)
		return
	}

	// Ensure at least one field is being updated
	if req.Name == nil && req.Email == nil && req.Phone == nil && req.UserType == nil {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}
	authUserID := middleware.GetUserID(r)
	err := db.UpdateUser(authUserID, userID, &req)

	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			http.Error(w, "Email or phone already exists", http.StatusConflict)
			return
		}
		log.Println("UpdateUser error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User updated successfully",
	})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")

	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "failed to begin transaction", http.StatusInternalServerError)
		return
	}
	defer db.TxFinalizer(tx, &err)

	// Check if user is assigned any assets
	assignedCount, err := db.NumberOfAssetsAssigned(tx, userID)
	if err != nil {
		http.Error(w, "failed to check user assignment", http.StatusInternalServerError)
		return
	}
	if assignedCount > 0 {
		http.Error(w, "user cannot be deleted while assets are assigned", http.StatusBadRequest)
		return
	}

	// Soft delete user
	err = db.SoftDeleteUser(tx, userID)
	if err != nil {
		http.Error(w, "failed to archive user", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("user deleted successfully"))
}
