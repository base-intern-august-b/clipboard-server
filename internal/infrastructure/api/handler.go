package api

import (
	"encoding/json"
	"net/http"

	"github.com/base-intern-august-b/clipboard-server/internal/service"
	"github.com/gorilla/mux"
)

// Handler はHTTPリクエストを処理するハンドラー
type Handler struct {
	router     *mux.Router
	appService *service.AppService
}

// NewHandler は新しいHandlerを作成する
func NewHandler(appService *service.AppService) *Handler {
	h := &Handler{
		router:     mux.NewRouter(),
		appService: appService,
	}

	// ルーティングの設定
	h.router.HandleFunc("/ping", h.pingHandler).Methods(http.MethodGet)

	return h
}

// ServeHTTP はhttp.Handlerインターフェースを実装する
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

// pingHandler は/pingエンドポイントを処理する
func (h *Handler) pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"message": "pong",
	})
}
