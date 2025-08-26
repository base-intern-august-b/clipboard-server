package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/base-intern-august-b/clipboard-server/internal/infrastructure/api"
	"github.com/base-intern-august-b/clipboard-server/internal/infrastructure/persistence"
	"github.com/base-intern-august-b/clipboard-server/internal/infrastructure/persistence/mysql"
	"github.com/base-intern-august-b/clipboard-server/internal/pkg/migration"
	"github.com/base-intern-august-b/clipboard-server/internal/usecase"

	"github.com/jmoiron/sqlx"
)

func main() {
	// 環境変数から設定を読み込む
	serverPort := getEnv("SERVER_PORT", "8080")

	log.Printf("Server starting on port %s", serverPort)

	// データベース接続
	db, err := sqlx.Connect("mysql", persistence.MySQL().FormatDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// マイグレーション実行
	log.Println("Running database migrations...")
	if err := migration.MigrateTables(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed successfully")

	// リポジトリの初期化
	userRepo := mysql.NewUserRepository(db)
	messageRepo := mysql.NewMessageRepository(db)
	channelRepo := mysql.NewChannelRepository(db)

	// ユースケースの初期化
	userUsecase := usecase.NewUserUsecase(userRepo)
	messageUsecase := usecase.NewMessageUsecase(messageRepo)
	channelUsecase := usecase.NewChannelUsecase(channelRepo)

	// APIルーターの設定
	router := api.NewRouter(channelUsecase, messageUsecase, userUsecase)
	handler := router.Setup()

	// HTTPサーバーの設定
	server := &http.Server{
		Addr:    ":" + serverPort,
		Handler: handler,
	}

	// グレースフルシャットダウンの設定
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// シグナル待機
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// グレースフルシャットダウン
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
