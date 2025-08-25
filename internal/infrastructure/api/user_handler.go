package api

import (
	"encoding/json"
	"net/http"

	"github.com/base-intern-august-b/clipboard-server/internal/domain/model"
	"github.com/base-intern-august-b/clipboard-server/internal/domain/usecase"
	"github.com/go-chi/chi/v5"
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

// GetUserByID : GET /v1/users/{userName}
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userName := chi.URLParam(r, "userName")

	user, err := h.userUsecase.GetUserByName(r.Context(), userName)
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

// GetUsersByID : POST /v1/users/_batch
func (h *UserHandler) GetUsersByID(w http.ResponseWriter, r *http.Request) {
	var req model.RequestGetUserBatch
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	users, err := h.userUsecase.GetUsersByName(r.Context(), req.UserNames)
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// PatchUser : PATCH /v1/users/{userName}
func (h *UserHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userName := chi.URLParam(r, "userName")
	if userName == "" {
		http.Error(w, "userName is required", http.StatusBadRequest)
		return
	}

	var req model.RequestPatchUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.userUsecase.PatchUser(ctx, userName, &req)
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

// DeleteUser : DELETE /v1/users/{userName}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userName := chi.URLParam(r, "userName")
	if userName == "" {
		http.Error(w, "userName is required", http.StatusBadRequest)
		return
	}

	err := h.userUsecase.DeleteUser(ctx, userName)
	if err != nil {
		http.Error(w, "failed to delete user", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
