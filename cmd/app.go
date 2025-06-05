package cmd

import (
	"1337bo4rd/internal/adapter/api"
	"1337bo4rd/internal/adapter/config"
	"1337bo4rd/internal/adapter/logger"
	"1337bo4rd/internal/adapter/storage/postgres"
	"1337bo4rd/internal/adapter/storage/minio"
	"1337bo4rd/internal/adapter/storage/postgres/repository"
	"1337bo4rd/internal/core/service"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	flag "1337bo4rd/internal/adapter/config"
	httpserver "1337bo4rd/internal/adapter/handler/http"
)

func Run() {
	// Parse flags
	minio, err := minio.NewMinioClient(
		"minio:9000", // endpoint
	)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize Minio client: %v", err))
	}

	err = flag.Parse()
	if err != nil {
		log.Fatal(err)
	}

	// Load environment variables
	config := config.New()

	// Set logger
	logger.Set(config.App)
	slog.Info("Staring application", "app", config.App.Name, "env", config.App.Env)

	// Init database
	db, err := postgres.OpenDB(config.DB)
	if err != nil {
		slog.Error("Failed to connect to the database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Avatar provider
	avatarProv := api.NewRickAndMorty()

	// Posts
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	userRepo := repository.NewUserRepository(db)
	// Services
	postService := service.NewPostService(postRepo, commentRepo)
	userService := service.NewUserService(userRepo, avatarProv)
	// Handlers
	postHandler := httpserver.NewPostHandler(postService, minio)

	mux := httpserver.NewRouter(*postHandler, userService)


	slog.Info(fmt.Sprintf("Listening on port: %d", flag.Port))
	err = http.ListenAndServe(fmt.Sprintf(":%d", flag.Port), mux)
	log.Fatal(err)
}
