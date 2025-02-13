package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/strahe/piecehub/api"
	"github.com/strahe/piecehub/config"
	"github.com/strahe/piecehub/storage"
	"github.com/strahe/piecehub/version"
	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:    "piecehub",
		Usage:   "A piece storage service",
		Version: version.GetVersion(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "config.toml",
				Usage:   "load configuration from `FILE`",
			},
			&cli.StringSliceFlag{
				Name:  "token",
				Usage: "token for accessing the service, can specify multiple tokens",
			},
			&cli.StringFlag{
				Name:    "log-level",
				Value:   "info",
				Usage:   "log level (debug, info, warn, error)",
				EnvVars: []string{"LOG_LEVEL"},
			},
		},
		Commands: []*cli.Command{
			dirCmd,
			s3Cmd,
		},
		Action: func(c *cli.Context) error {
			configPath := c.String("config")
			if !filepath.IsAbs(configPath) {
				pwd, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("get working directory: %v", err)
				}
				configPath = filepath.Join(pwd, configPath)
			}

			cfg, err := config.LoadConfig(configPath)
			if err != nil {
				return fmt.Errorf("load config: %v", err)
			}

			tokens := c.StringSlice("token")
			if len(tokens) > 0 {
				cfg.Server.Tokens = append(cfg.Server.Tokens, tokens...)
			}

			return runServer(cfg)
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func runServer(cfg *config.Config) error {
	store, err := storage.NewManager(cfg)
	if err != nil {
		return fmt.Errorf("create storage manager: %v", err)
	}

	handler := api.NewHandler(cfg, store)

	log.Printf("Starting server on %s", cfg.Server.Address)
	if err := startServer(cfg, handler); err != nil {
		return fmt.Errorf("start server: %v", err)
	}

	return nil
}

func startServer(cfg *config.Config, handler http.Handler) error {
	server := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      handler,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	return server.ListenAndServe()
}
