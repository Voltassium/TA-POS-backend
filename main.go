package main

import (
	"backend-ta/app/database/seeders"
	"backend-ta/app/routes"
	"backend-ta/pkg/authentication"
	"backend-ta/pkg/config"
	"backend-ta/pkg/database"
	"backend-ta/pkg/http/server"
	"backend-ta/pkg/logger"
	"flag"
	"log"
)

func main() {
	cfg := config.LoadConfig()

	logger.NewZapLogger(cfg.Logger)

	database.InitDB(cfg.Database)

	authentication.NewJWTManager(authentication.JWTOptions{
		AccessSecret:       cfg.Authentication.AccessSecretKey,
		RefreshSecret:      cfg.Authentication.RefreshSecretKey,
		Issuer:             cfg.Authentication.Issuer,
		ExpiryAccessToken:  cfg.Authentication.AccessTokenExpiry,
		ExpiryRefreshToken: cfg.Authentication.RefreshTokenExpiry,
	})

	authentication.SetupKey(cfg.Authentication.EncryptKey)

	seed := flag.Bool("seed", false, "seed database")
	fresh := flag.Bool("fresh", false, "truncate tables before seeding (use with --seed)")
	flag.Parse()

	if *seed {
		err := seeders.SeedAll(database.GetDB().DB, *fresh)
		if err != nil {
			log.Fatal("Failed to seed database", err)
		}
		return
	}

	server.Init(cfg.Application, routes.RegisterV1).GracefulShutdown()
	defer func() {
		err := logger.Log.Sync()
		if err != nil {
			log.Println(err)
		}
	}()
}
