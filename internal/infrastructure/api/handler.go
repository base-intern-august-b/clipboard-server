package api

import (
	"fmt"
	"net/http"

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
