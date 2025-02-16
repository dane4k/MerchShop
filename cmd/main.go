package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dane4k/MerchShop/internal/infrastructure/config"
	"github.com/dane4k/MerchShop/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.MustLoad()
	if err != nil {
		logrus.Fatal(err)
	}

	logger.InitLogger(cfg)

	appObj := initApp(cfg)
	defer appObj.Close()

	router := gin.Default()
	initializeRoutes(router, appObj)

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
