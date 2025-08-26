package api

import (
	"encoding/json"
	"net/http"

	"github.com/base-intern-august-b/clipboard-server/internal/domain/model"
	"github.com/base-intern-august-b/clipboard-server/internal/domain/usecase"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

// CreateUser : POST /v1/users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req model.RequestCreateUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userUsecase.CreateUser(r.Context(), &req)
	if err != nil {
		if err == model.ErrAlreadyExistUserName {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUsers : GET /v1/users
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userUsecase.GetUsers(r.Context())
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetUserByID : GET /v1/users/{userID}
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := getID(r, "userID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.userUsecase.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// PatchUser : PATCH /v1/users/{userId}
func (h *UserHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	userID, err := getID(r, "userID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req model.RequestPatchUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.userUsecase.PatchUser(r.Context(), userID, &req)
	if err != nil {
		if err == model.ErrAlreadyExistUserName {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "failed to patch user", http.StatusBadRequest)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// ChangePassword : POST /v1/users/{userID}/change-password
func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, err := getID(r, "userID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req model.RequestChangePassword
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err = h.userUsecase.ChangePassword(r.Context(), userID, &req)
	if err != nil {
		if err == model.ErrNothingChanged {
			http.Error(w, err.Error(), http.StatusNoContent)
			return
		}
		http.Error(w, "failed to change password", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteUser : DELETE /v1/users/{userID}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := getID(r, "userID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.userUsecase.DeleteUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to delete user", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
