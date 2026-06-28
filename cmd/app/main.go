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
	cfg := config.Load()

	conn := database.InitDB(cfg.DatabaseURL)
	defer conn.Close(context.Background())

	userRepo := repository.NewRepo(conn)

	
	authService := service.NewAuthService(
		userRepo, 
		cfg.SessionDuration, 
		cfg.LockoutDuration,
	)

	cliHandler := cli.NewCLIHandler(authService, userRepo)

	cliHandler.Run()
}