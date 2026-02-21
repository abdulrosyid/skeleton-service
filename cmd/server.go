package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	cfg "skeleton-service/configs"
	"skeleton-service/internal/adapter/inbound/registry"
	"skeleton-service/internal/adapter/inbound/rest"
	"skeleton-service/internal/adapter/outbound/datasource"
	repoRegistry "skeleton-service/internal/adapter/outbound/repository/mysqldb"
)

func RunServer() {
	configuration := cfg.GetConfig()

	logrus.SetFormatter(&logrus.JSONFormatter{})

	db, err := datasource.NewMySqlDB(configuration.MySQLDSN())
	if err != nil {
		log.Fatalf("unable to initialize mysql: %v", err)
	}
	defer db.Close()

	repositories := repoRegistry.NewRepoSQL(db)
	serviceRegistry := registry.NewServiceRegistry(repositories, configuration)

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	basicAuthGroup := e.Group("")

	rest.Apply(e, basicAuthGroup, *serviceRegistry)

	go func() {
		if err := e.Start(fmt.Sprintf(":%s", configuration.HTTPPort)); err != nil {
			e.Logger.Fatalf("shutting down the server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ServerShutdownTimeout())
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatalf("server forced to shutdown: %v", err)
	}
}
