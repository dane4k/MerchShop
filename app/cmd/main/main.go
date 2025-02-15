package main

import (
	"MerchShop/app"
	"MerchShop/internal/config"
	"MerchShop/internal/logger"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// TODO: тесты, линтер, рефакторинг

func main() {
	cfg, err := config.MustLoad()
	if err != nil {
		logrus.Fatal(err)
	}

	logger.InitLogger(cfg)

	appObj := app.InitApp(cfg)
	defer appObj.Close()

	router := gin.Default()
	app.InitializeRoutes(router, appObj)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.Server.Port),
		Handler: router,
	}

	go func() {
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	logrus.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		logrus.WithError(err).Fatal("server forced to shutdown:")
	}

	logrus.Info("server exited")

}
