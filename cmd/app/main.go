package main

import (
	"database/sql"
	"log"
	"mfawzanid/warehouse-commerce/core/repository"
	"mfawzanid/warehouse-commerce/core/usecase"
	"mfawzanid/warehouse-commerce/generated"
	"mfawzanid/warehouse-commerce/handler"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/lib/pq"
)

func main() {
	e := echo.New()

	dbDSN := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbDSN)
	if err != nil {
		panic(err)
	}
	log.Printf("Connected to database: %v", dbDSN)
	defer db.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Protocol: 2, // connection protocol
	})

	// repository
	inventoryRepo := repository.NewInventoryRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	redisRepo := repository.NewRedisRepository(redisClient)
	userRepo := repository.NewUserRepository(db)

	// usecase
	authUsecase := usecase.NewAuthUsecase()
	userUsecase := usecase.NewUserUsecase(userRepo, authUsecase)
	inventoryUsecase := usecase.NewInventoryUsecase(inventoryRepo)
	transactionUsecase := usecase.NewTransactionUsecase(inventoryRepo, transactionRepo, redisRepo)

	// handler
	authHandler := handler.NewAuthHandler(authUsecase)
	serverHandler := handler.NewServer(userUsecase, inventoryUsecase, transactionUsecase)
	var server generated.ServerInterface = serverHandler

	// protected routes
	protectedGroup := e.Group("")
	protectedGroup.Use(authHandler.VerifyToken)
	generated.RegisterHandlers(protectedGroup, server)

	// public routes
	e.GET("/health", serverHandler.GetHealth)
	e.POST("/user/register", serverHandler.RegisterUser)
	e.POST("/user/login", serverHandler.Login)

	e.Logger.Fatal(e.Start(":1323"))
}
