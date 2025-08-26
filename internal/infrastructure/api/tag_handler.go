package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/base-intern-august-b/clipboard-server/internal/domain/usecase"
	"github.com/go-chi/chi/v5"
)

type TagHandler struct {
	tagUsecase usecase.TagUsecase
}

func NewTagHandler(uc usecase.TagUsecase) *TagHandler {
	return &TagHandler{tagUsecase: uc}
}

// ModifyMessageTags : PUT /v1/messages/{messageID}/tags
func (h *TagHandler) ModifyMessageTags(w http.ResponseWriter, r *http.Request) {
	messageID, err := getID(r, "messageID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req struct {
		Tags []string `json:"tags"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	messageDetail, err := h.tagUsecase.ModifyMessageTags(r.Context(), messageID, req.Tags)
	if err != nil {
		// Add more specific error handling based on usecase errors
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messageDetail)
}

// FindMessagesByTag : GET /v1/tags/{tagName}/messages
func (h *TagHandler) FindMessagesByTag(w http.ResponseWriter, r *http.Request) {
	tagName := chi.URLParam(r, "tagName")
	limit, offset := getLimitOffset(r)

	messageDetails, err := h.tagUsecase.FindMessagesByTag(r.Context(), tagName, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messageDetails)
}

// FindMessagesByTagInChannel : GET /v1/channels/{channelID}/tags/{tagName}/messages
func (h *TagHandler) FindMessagesByTagInChannel(w http.ResponseWriter, r *http.Request) {
	channelID, err := getID(r, "channelID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tagName := chi.URLParam(r, "tagName")
	limit, offset := getLimitOffset(r)

	messageDetails, err := h.tagUsecase.FindMessagesByTagInChannel(r.Context(), channelID, tagName, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messageDetails)
}
