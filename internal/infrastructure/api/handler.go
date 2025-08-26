package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/base-intern-august-b/clipboard-server/internal/domain/model"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
)

func getID(r *http.Request, key string) (uuid.UUID, error) {
	id, err := uuid.FromString(chi.URLParam(r, key))
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %s", model.ErrInvalidUUID, err.Error())
	} else if id.IsNil() {
		return uuid.Nil, model.ErrNilUUID
	}

	return id, nil
}

func getLimitOffset(r *http.Request) (int, int) {
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limitStr == "" {
		limit = 100 // Default limit
	}

	offsetStr := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offsetStr == "" {
		offset = 0 // Default offset
	}
	return limit, offset
}
