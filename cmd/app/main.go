package main

import (
	"context"

	"github.com/aravindgopan-dev/cli-auth-system/internal/cli"
	"github.com/aravindgopan-dev/cli-auth-system/internal/config"
	"github.com/aravindgopan-dev/cli-auth-system/internal/database"
	"github.com/aravindgopan-dev/cli-auth-system/internal/repository"
	"github.com/aravindgopan-dev/cli-auth-system/internal/service"
)

func main() {
	// 1. Initialize configuration values centrally
	cfg := config.Load()

	// 2. Open raw connection using the pure pgx library
	conn := database.InitDB(cfg.DatabaseURL)
	defer conn.Close(context.Background())

	// 3. Create the concrete repository
	userRepo := repository.NewRepo(conn)

	
	authService := service.NewAuthService(
		userRepo, 
		cfg.SessionDuration, 
		cfg.LockoutDuration,
	)

	// 5. Inject the service and repository into the front-facing CLI prompt loop.
	// userRepo implicitly satisfies cli.SessionStore perfectly.
	cliHandler := cli.NewCLIHandler(authService, userRepo)

	cliHandler.Run()
}