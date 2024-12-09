package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/kholidss/xyz-skilltest/internal/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
	"github.com/kholidss/xyz-skilltest/internal/bootstrap"
	"github.com/kholidss/xyz-skilltest/pkg/config"
)

type App struct {
	*fiber.App
	Cfg *config.Config
}

var appServer *App

func InitializeApp(cfg *config.Config) {

	// boostrap run and initialize package dependency
	bootstrap.RegistryLogger(cfg)
	otelManager, err := bootstrap.RegistryOpenTelemetry(cfg)
	if err != nil {
		otelManager.Close()
	}

	f := fiber.New(cfg.FiberConfig())
	f.Use(
		cors.New(cors.Config{
			MaxAge: 300,
			AllowOrigins: strings.Join([]string{
				"http://*",
				"https://*",
			}, ","),
			AllowHeaders: strings.Join([]string{
				"Origin",
				"Content-Type",
				"Accept",
			}, ","),
			AllowMethods: strings.Join([]string{
				fiber.MethodGet,
				fiber.MethodPost,
				fiber.MethodPut,
				fiber.MethodDelete,
				fiber.MethodHead,
			}, ","),
		}),
		requestid.New(requestid.Config{
			ContextKey: "request-id",
			Header:     "X-Request-ID",
			Generator: func() string {
				return uuid.New().String()
			},
		}),
		logger.New(logger.Config{
			Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
		}),
	)

	rtr := router.NewRouter(cfg, f)
	rtr.Route()

	appServer = &App{
		App: f,
		Cfg: cfg,
	}
}

func (app *App) StartServer() error {
	// Run the server in a separate goroutine
	go func() {
		if err := app.Listen(fmt.Sprintf("%v:%v", app.Cfg.AppHost, app.Cfg.AppPort)); err != nil {
			log.Fatalf("Failed starting server: %v\n", err)
		}
	}()

	// Channel to listen for system signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Wait for a signal
	<-quit

	// Graceful shutdown
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("Error during server shutdown: %v\n", err)
		return err
	}

	log.Println("Server gracefully stopped")
	return nil
}

func GetServer() *App {
	return appServer
}
