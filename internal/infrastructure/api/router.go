package api

import (
	"net/http"

	"github.com/base-intern-august-b/clipboard-server/internal/domain/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Router struct {
	channelUsecase usecase.ChannelUsecase
	messageUsecase usecase.MessageUsecase
	userUsecase    usecase.UserUsecase
}

func NewRouter(channelUsecase usecase.ChannelUsecase, messageUsecase usecase.MessageUsecase, userUsecase usecase.UserUsecase) *Router {
	return &Router{
		channelUsecase: channelUsecase,
		messageUsecase: messageUsecase,
		userUsecase:    userUsecase,
	}
}

func (r *Router) Setup() http.Handler {
	router := chi.NewRouter()

	// ミドルウェアの設定
	router.Use(LoggingMiddleware)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Route("/api/v1", func(v1 chi.Router) {
		// ユーザーAPI
		userHandler := NewUserHandler(r.userUsecase)
		v1.Route("/users", func(user chi.Router) {
			user.Post("/", userHandler.CreateUser)
			user.Get("/", userHandler.GetUsers)
			user.Get("/{userID}", userHandler.GetUserByID)
			user.Patch("/{userID}", userHandler.PatchUser)
			user.Post("/{userID}/change-password", userHandler.ChangePassword)
			user.Delete("/{userID}", userHandler.DeleteUser)
		})

		// チャンネルAPI
		channelHandler := NewChannelHandler(r.channelUsecase)
		messageHandler := NewMessageHandler(r.messageUsecase)
		v1.Route("/channels", func(channel chi.Router) {
			channel.Post("/", channelHandler.CreateChannel)
			channel.Get("/", channelHandler.GetChannels)

			channel.Route("/{channelID}", func(ch chi.Router) {
				ch.Get("/", channelHandler.GetChannelByName)
				ch.Patch("/", channelHandler.PatchChannel)
				ch.Delete("/", channelHandler.DeleteChannel)

				// チャンネルごとのメッセージ
				ch.Get("/messages", messageHandler.GetMessages)
				ch.Get("/messages/span", messageHandler.GetMessagesInDuration)
				ch.Get("/messages/pinned", messageHandler.GetPinnedMessages)
			})
		})

		// メッセージAPI
		v1.Route("/messages", func(message chi.Router) {
			message.Post("/", messageHandler.CreateMessage)
			message.Patch("/{messageID}", messageHandler.PatchMessage)
			message.Delete("/{messageID}", messageHandler.DeleteMessage)
			message.Post("/{messageID}/pin", messageHandler.PinnMessage)
			message.Post("/{messageID}/unpin", messageHandler.UnpinnMessage)
		})
	})

	// 静的ファイルの配信（CSS、JS、画像など）
	router.Get("/styles.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./view/styles.css")
	})
	router.Get("/script.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./view/script.js")
	})

	// ルートページ
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./view/index.html")
	})

	// その他の静的ファイル（フォールバック）
	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		// APIエンドポイント以外のリクエストは静的ファイルとして処理
		filePath := "./view" + r.URL.Path
		http.ServeFile(w, r, filePath)
	})

	return router
}
