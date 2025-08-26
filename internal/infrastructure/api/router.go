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
	tagUsecase     usecase.TagUsecase
}

func NewRouter(channelUsecase usecase.ChannelUsecase, messageUsecase usecase.MessageUsecase, userUsecase usecase.UserUsecase, tagUsecase usecase.TagUsecase) *Router {
	return &Router{
		channelUsecase: channelUsecase,
		messageUsecase: messageUsecase,
		userUsecase:    userUsecase,
		tagUsecase:     tagUsecase,
	}
}

func (r *Router) Setup() http.Handler {
	router := chi.NewRouter()

	// Middleware
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

	// Handlers
	userHandler := NewUserHandler(r.userUsecase)
	channelHandler := NewChannelHandler(r.channelUsecase)
	messageHandler := NewMessageHandler(r.messageUsecase)
	tagHandler := NewTagHandler(r.tagUsecase)

	router.Route("/api/v1", func(v1 chi.Router) {
		// User API
		v1.Route("/users", func(user chi.Router) {
			user.Post("/", userHandler.CreateUser)
			user.Get("/", userHandler.GetUsers)
			user.Get("/{userID}", userHandler.GetUserByID)
			user.Patch("/{userID}", userHandler.PatchUser)
			user.Post("/{userID}/change-password", userHandler.ChangePassword)
			user.Delete("/{userID}", userHandler.DeleteUser)
		})

		// Channel API
		v1.Route("/channels", func(channel chi.Router) {
			channel.Post("/", channelHandler.CreateChannel)
			channel.Get("/", channelHandler.GetChannels)

			channel.Route("/{channelID}", func(ch chi.Router) {
				ch.Get("/", channelHandler.GetChannelByName)
				ch.Patch("/", channelHandler.PatchChannel)
				ch.Delete("/", channelHandler.DeleteChannel)

				// Messages in Channel
				ch.Get("/messages", messageHandler.GetMessages)
				ch.Get("/messages/span", messageHandler.GetMessagesInDuration)
				ch.Get("/messages/pinned", messageHandler.GetPinnedMessages)

				// Tags in Channel
				ch.Get("/tags/{tagName}/messages", tagHandler.FindMessagesByTagInChannel)
			})
		})

		// Message API
		v1.Route("/messages", func(message chi.Router) {
			message.Post("/", messageHandler.CreateMessage)
			message.Get("/{messageID}", messageHandler.GetMessage)
			message.Patch("/{messageID}", messageHandler.PatchMessage)
			message.Delete("/{messageID}", messageHandler.DeleteMessage)
			message.Post("/{messageID}/pin", messageHandler.PinnMessage)
			message.Post("/{messageID}/unpin", messageHandler.UnpinnMessage)
			// Tag modification for a message
			message.Put("/{messageID}/tags", tagHandler.ModifyMessageTags)
		})

		// Tag API
		v1.Route("/tags", func(tag chi.Router) {
			tag.Get("/{tagName}/messages", tagHandler.FindMessagesByTag)
		})
	})

	return router
}
