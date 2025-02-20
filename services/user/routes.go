package user

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/varnit-ta/Ecom-API/configs"
	"github.com/varnit-ta/Ecom-API/services/auth"
	"github.com/varnit-ta/Ecom-API/types"
	"github.com/varnit-ta/Ecom-API/utils"
)

/*
Handler provides HTTP request handlers related to user authentication and management.

@param store - types.UserStore: The user store to interact with the database.

@return *Handler: A new instance of Handler.
*/
type Handler struct {
	store types.UserStore
}

/*
NewHandler creates a new user handler with the provided store.

@param store - types.UserStore: The store used for interacting with the user database.

@return *Handler: A new Handler instance.
*/
func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

/*
RegisterRoutes registers user-related routes to the provided router.

@param router - *mux.Router: The router to register user routes on.

@return void
*/
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

/*
handleLogin handles user authentication by verifying credentials and returning a JWT token.

@param w - http.ResponseWriter: The response writer to send HTTP responses.
@param r - *http.Request: The incoming HTTP request.

@return void
*/
func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var user types.LoginUserPayload

	// Parse JSON request body into user struct
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Validate request payload
	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	// Retrieve user by email
	u, err := h.store.GetUserByEmail(user.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	// Compare provided password with hashed password
	if !auth.ComparePasswords(u.Password, []byte(user.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	// Generate JWT token for authentication
	secret := []byte(configs.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Return token in response
	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

/*
handleRegister handles user registration by creating a new user in the database.

@param w - http.ResponseWriter: The response writer to send HTTP responses.
@param r - *http.Request: The incoming HTTP request.

@return void
*/
func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var user types.RegisterUserPayload

	// Parse JSON request body into user struct
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Validate request payload
	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	// Check if user already exists
	_, err := h.store.GetUserByEmail(user.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", user.Email))
		return
	}

	// Hash password before storing it
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Create new user record
	err = h.store.CreateUser(types.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Return success response
	utils.WriteJSON(w, http.StatusCreated, nil)
}
