package api

import (
	"encoding/json"
	"net/http"

	"github.com/base-intern-august-b/clipboard-server/internal/domain/model"
	"github.com/base-intern-august-b/clipboard-server/internal/domain/usecase"
)

type MessageHandler struct {
	messageUsecase usecase.MessageUsecase
}

func NewMessageHandler(messageUsecase usecase.MessageUsecase) *MessageHandler {
	return &MessageHandler{messageUsecase: messageUsecase}
}

// CreateMessage : POST /v1/messages
func (h *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var req model.RequestCreateMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	message, err := h.messageUsecase.CreateMessage(r.Context(), &req)
	if err != nil {
		if err == model.ErrChannelNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else if err == model.ErrInvalidMessageContent {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

// GetMessages : GET /v1/messages/channel
func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	var req model.RequestGetMessages
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	messages, err := h.messageUsecase.GetMessages(r.Context(), &req)
	if err != nil {
		if err == model.ErrChannelNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// GetMessagesInDuration : POST /v1/messages/span
func (h *MessageHandler) GetMessagesInDuration(w http.ResponseWriter, r *http.Request) {
	var req model.RequestGetMessagesInDuration
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	messages, err := h.messageUsecase.GetMessagesInDuration(r.Context(), &req)
	if err != nil {
		if err == model.ErrChannelNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else if err == model.ErrInvalidTimeRange {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// GetPinnedMessages : GET /v1/messages/pinned/{channelID}
func (h *MessageHandler) GetPinnedMessages(w http.ResponseWriter, r *http.Request) {
	channelID, err := getID(r, "channelID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	messages, err := h.messageUsecase.GetPinnedMessages(r.Context(), channelID)
	if err != nil {
		if err == model.ErrChannelNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// PatchMessage : PATCH /v1/messages/{messageID}
func (h *MessageHandler) PatchMessage(w http.ResponseWriter, r *http.Request) {
	messageID, err := getID(r, "messageID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req model.RequestPatchMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	message, err := h.messageUsecase.PatchMessage(r.Context(), messageID, &req)
	if err != nil {
		if err == model.ErrMessageNotFound || err == model.ErrChannelNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else if err == model.ErrInvalidMessageContent {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

// PinnMessage : POST /v1/messages/pinn/{messageID}
func (h *MessageHandler) PinnMessage(w http.ResponseWriter, r *http.Request) {
	messageID, err := getID(r, "messageID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.messageUsecase.PinnMessage(r.Context(), messageID)
	if err != nil {
		if err == model.ErrMessageNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UnpinnMessage : POST /v1/messages/unpinn/{messageID}
func (h *MessageHandler) UnpinnMessage(w http.ResponseWriter, r *http.Request) {
	messageID, err := getID(r, "messageID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.messageUsecase.UnpinnMessage(r.Context(), messageID)
	if err != nil {
		if err == model.ErrMessageNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteMessage : DELETE /v1/messages/{messageID}
func (h *MessageHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	messageID, err := getID(r, "messageID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.messageUsecase.DeleteMessage(r.Context(), messageID)
	if err != nil {
		if err == model.ErrMessageNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
