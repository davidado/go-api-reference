// Package user : User service
package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/davidado/go-api-reference/config"
	"github.com/davidado/go-api-reference/netjson"
	"github.com/davidado/go-api-reference/service/auth"
	"github.com/davidado/go-api-reference/types"
	vd "github.com/davidado/go-api-reference/validator"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

// Handler : User handler
type Handler struct {
	store types.UserStore
}

// NewHandler : Create a new user handler
func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

// RegisterRoutes : Register user routes
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods(http.MethodPost)
	router.HandleFunc("/register", h.handleRegister).Methods(http.MethodPost)

	// admin route
	router.HandleFunc("/users/{userID}", auth.WithJWTAuth(h.handleGetUser, h.store)).Methods(http.MethodGet)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// get JSON payload.
	var payload types.LoginUserPayload
	if err := netjson.Parse(r, &payload); err != nil {
		netjson.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate the payload.
	if err := vd.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		netjson.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		netjson.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	if !auth.ComparePasswords(u.Password, []byte(payload.Password)) {
		netjson.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		netjson.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	netjson.Write(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// get JSON payload.
	var payload types.RegisterUserPayload
	if err := netjson.Parse(r, &payload); err != nil {
		netjson.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate the payload.
	if err := vd.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		netjson.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// Check if the user exists.
	_, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		netjson.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		netjson.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// if it doesn't, we create the new user.
	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		netjson.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	netjson.Write(w, http.StatusCreated, map[string]string{"message": "user created"})
}

func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["userID"]
	if !ok {
		netjson.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user ID"))
		return
	}

	userID, err := strconv.Atoi(str)
	if err != nil {
		netjson.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}

	user, err := h.store.GetUserByID(userID)
	if err != nil {
		netjson.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	netjson.Write(w, http.StatusOK, user)
}
