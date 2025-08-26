package api

import (
	"encoding/json"
	"net/http"

	"github.com/base-intern-august-b/clipboard-server/internal/domain/model"
	"github.com/base-intern-august-b/clipboard-server/internal/domain/usecase"
	"github.com/go-chi/chi/v5"
)

type ChannelHandler struct {
	channelUsecase usecase.ChannelUsecase
}

func NewChannelHandler(channelUsecase usecase.ChannelUsecase) *ChannelHandler {
	return &ChannelHandler{channelUsecase: channelUsecase}
}

// CreateChannel : POST /v1/channels
func (h *ChannelHandler) CreateChannel(w http.ResponseWriter, r *http.Request) {
	var req model.RequestCreateChannel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	channel, err := h.channelUsecase.CreateChannel(r.Context(), &req)
	if err != nil {
		if err == model.ErrAlreadyExistChannelName {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(channel)
}

// GetChannelByName : GET /v1/channels/{channelName}
func (h *ChannelHandler) GetChannelByName(w http.ResponseWriter, r *http.Request) {
	channelName := chi.URLParam(r, "channelName")

	channel, err := h.channelUsecase.GetChannelByName(r.Context(), channelName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if channel == nil {
		http.Error(w, "Channel not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(channel)
}

// GetChannels : POST /v1/channels
func (h *ChannelHandler) GetChannels(w http.ResponseWriter, r *http.Request) {
	channels, err := h.channelUsecase.GetChannels(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(channels)
}

// PatchChannel : PATCH /v1/channels/{channelName}
func (h *ChannelHandler) PatchChannel(w http.ResponseWriter, r *http.Request) {
	channelName := chi.URLParam(r, "channelName")

	var req model.RequestPatchChannel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	channel, err := h.channelUsecase.PatchChannel(r.Context(), channelName, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if channel == nil {
		http.Error(w, "Channel not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(channel)
}

func (h *ChannelHandler) DeleteChannel(w http.ResponseWriter, r *http.Request) {
	channelName := chi.URLParam(r, "channelName")

	err := h.channelUsecase.DeleteChannel(r.Context(), channelName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
