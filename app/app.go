package app

import (
	"MerchShop/internal/config"
	"MerchShop/internal/db"
	"MerchShop/internal/handler"
	"MerchShop/internal/middleware"
	"MerchShop/internal/repo/pgdb"
	"MerchShop/internal/service"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type App struct {
	DB                 *sql.DB
	InventoryRepo      pgdb.InventoryRepo
	TransactionRepo    pgdb.TransactionRepo
	UserRepo           pgdb.UserRepo
	TransactionService service.TransactionService
	AuthService        service.AuthService
	UserService        service.UserService
	AuthHandler        *handler.AuthHandler
	InfoHandler        *handler.InfoHandler
	InventoryHandler   *handler.InventoryHandler
	TransactionHandler *handler.TransactionHandler
	AuthMiddleware     *middleware.AuthMiddleware
}

func InitApp(cfg *config.Config) *App {
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

	authHandler := handler.NewAuthHandler(userService, baseHandler)
	infoHandler := handler.NewInfoHandler(userService)
	inventoryHandler := handler.NewInventoryHandler(userService, baseHandler)
	transactionHandler := handler.NewTransactionHandler(userService, baseHandler)

	authMiddleware := middleware.NewAuthMiddleware(authService)

	return &App{
		DB:                 DB,
		InventoryRepo:      inventoryRepo,
		TransactionRepo:    transactionRepo,
		UserRepo:           userRepo,
		TransactionService: transactionService,
		AuthService:        authService,
		UserService:        userService,
		AuthHandler:        authHandler,
		InfoHandler:        infoHandler,
		InventoryHandler:   inventoryHandler,
		TransactionHandler: transactionHandler,
		AuthMiddleware:     authMiddleware,
	}
}

func InitializeRoutes(router *gin.Engine, appObj *App) {
	api := router.Group("/api")
	api.Use(appObj.AuthMiddleware.Handle())
	{
		api.GET("/info", appObj.InfoHandler.GetInfo)
		api.POST("/sendCoin", appObj.TransactionHandler.SendCoins)
		api.GET("/buy/:item", appObj.InventoryHandler.BuyItem)
	}

	router.POST("/api/auth", appObj.AuthHandler.Auth)
}

func (appObj *App) Close() {
	if appObj.DB != nil {
		if err := appObj.DB.Close(); err != nil {
			logrus.Fatal(err)
		} else {
			logrus.Info("db connection closed")
		}
	}
}
