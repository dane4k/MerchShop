package main

import (
	"database/sql"

	"github.com/dane4k/MerchShop/internal/application/service"
	"github.com/dane4k/MerchShop/internal/infrastructure/config"
	"github.com/dane4k/MerchShop/internal/infrastructure/db"
	"github.com/dane4k/MerchShop/internal/infrastructure/db/pgdb"
	"github.com/dane4k/MerchShop/internal/presentations/handler"
	"github.com/dane4k/MerchShop/internal/presentations/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type App struct {
	db                 *sql.DB
	inventoryRepo      *pgdb.InventoryRepo
	transactionRepo    pgdb.TransactionRepo
	userRepo           pgdb.UserRepo
	transactionService service.TransactionService
	authService        service.AuthService
	userService        service.UserService
	authHandler        *handler.AuthHandler
	infoHandler        *handler.InfoHandler
	inventoryHandler   *handler.InventoryHandler
	transactionHandler *handler.TransactionHandler
	authMiddleware     *middleware.AuthMiddleware
}

func initApp(cfg *config.Config) *App {
	DB, err := db.InitDB(cfg)
	if err != nil {
		logrus.Fatal(err)
	}

	inventoryRepo := pgdb.NewInventoryRepo(DB)
	transactionRepo := pgdb.NewTransactionRepo(DB)
	userRepo := pgdb.NewUserRepo(DB, transactionRepo, inventoryRepo)

	transactionService := service.NewTransactionService(userRepo, transactionRepo)
	authService := service.NewAuthService(cfg.JWT.JWTSecret)
	userService := service.NewUserService(userRepo, inventoryRepo, transactionService, authService)

	baseHandler := handler.NewBaseHandler()

	return &App{
		db:                 DB,
		inventoryRepo:      inventoryRepo,
		transactionRepo:    *transactionRepo,
		userRepo:           *userRepo,
		transactionService: transactionService,
		authService:        authService,
		userService:        userService,
		authHandler:        handler.NewAuthHandler(userService, baseHandler),
		infoHandler:        handler.NewInfoHandler(userService),
		inventoryHandler:   handler.NewInventoryHandler(userService, baseHandler),
		transactionHandler: handler.NewTransactionHandler(userService, baseHandler),
		authMiddleware:     middleware.NewAuthMiddleware(authService),
	}
}

func initializeRoutes(router *gin.Engine, app *App) {
	api := router.Group("/api")
	api.Use(app.authMiddleware.Handle())
	{
		api.GET("/info", app.infoHandler.GetInfo)
		api.POST("/sendCoin", app.transactionHandler.SendCoins)
		api.GET("/buy/:item", app.inventoryHandler.BuyItem)
	}

	router.POST("/api/auth", app.authHandler.Auth)
}

func (appObj *App) Close() {
	if appObj.db != nil {
		if err := appObj.db.Close(); err != nil {
			logrus.Fatal(err)
		} else {
			logrus.Info("db connection closed")
		}
	}
}
