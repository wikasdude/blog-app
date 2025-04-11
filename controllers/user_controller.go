package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/wikasdude/blog-backend/model"
	service "github.com/wikasdude/blog-backend/services"
	"github.com/wikasdude/blog-backend/utils"
)

type UserController struct {
	userService service.UserService
}

// DTO object
type RegisterUserInput struct {
	Name     string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService: userService}
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Registers a new user with the provided name, email, and password.
// @Tags users
// @Accept  json
// @Produce  json
// @Param input body RegisterUserInput true "User registration details"
// @Success 201 {object} model.User
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 400 {object} map[string]string "All fields (username, email, password) are required"
// @Failure 400 {object} map[string]string "Invalid email format"
// @Failure 500 {object} map[string]string "Error hashing password"
// @Failure 500 {object} map[string]string "Failed to create user"
// @Router /users [post]
func (c *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var input RegisterUserInput
	//var user model.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if input.Name == "" || input.Email == "" || input.Password == "" {
		http.Error(w, "All fields (username, email, password) are required", http.StatusBadRequest)
		return
	}
	if !utils.IsValidEmail(input.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	role := "user" // Default role is "user"
	if input.Role == "admin" {
		role = "admin" // If the role provided in the input is "admin", set the role to "admin"
	}
	user := model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
		Role:     role,
	}
	fmt.Println("User to create:", user)
	err = c.userService.RegisterUser(&user)
	if err != nil {
		fmt.Print(err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

type LoginUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginUser godoc
// @Summary Login a user and generate JWT token
// @Description Logs in a user by verifying email and password, and returns a JWT token upon successful login.
// @Tags users
// @Accept  json
// @Produce  json
// @Param input body LoginUserInput true "User login details"
// @Success 200 {object} map[string]string "Login successful with JWT token"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 401 {object} map[string]string "Invalid email or password"
// @Failure 500 {object} map[string]string "Could not generate token"
// @Router /login [post]
func (c *UserController) LoginUser(w http.ResponseWriter, r *http.Request) {
	var input RegisterUserInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Fetch user from DB by email
	dbUser, err := c.userService.GetUserByEmail(input.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Compare the input password with the hashed password in DB
	if !utils.ComparePassword(dbUser.Password, input.Password) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(dbUser.ID, dbUser.Role)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token, "message": "login successful"})
}

// GetUserByID godoc
// @Summary Get user details by ID
// @Description Fetches a user by their unique ID.
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID" example("123")
// @Success 200 {object} model.User "User found" example(model.User{"id": "123", "name": "John Doe", "email": "johndoe@example.com"})
// @Failure 404 {object} map[string]string "User not found" example(map[string]string{"error": "User not found"})
// @Failure 400 {object} map[string]string "Invalid ID" example(map[string]string{"error": "Invalid ID"})
// @Router /users/{id} [get]
func (uc *UserController) GetUserByID(w http.ResponseWriter, r *http.Request, id string) {
	user, err := uc.userService.GetUserByID(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request, idStr string) {

	tokenString := r.Header.Get("Authorization")
	claims, err := utils.ValidateJWT(strings.TrimPrefix(tokenString, "Bearer "))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := claims.UserID
	role := claims.Role
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if role != "admin" && userID != id {
		http.Error(w, "Forbidden: You are not allowed to update this user", http.StatusForbidden)
		return
	}

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user.ID = id

	if err := c.userService.UpdateUser(&user); err != nil {
		fmt.Println("line no 63:", err)
		http.Error(w, fmt.Sprintf("Failed to update user %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = c.userService.DeleteUser(id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted successfully"))
}
