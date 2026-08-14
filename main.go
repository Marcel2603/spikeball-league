package main

import (
	"embed"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/Marcel2603/spikeball-league/cmd/config"
	"github.com/Marcel2603/spikeball-league/internal/db"
	"github.com/Marcel2603/spikeball-league/internal/handler/admin"
	"github.com/Marcel2603/spikeball-league/internal/handler/health"
	"github.com/Marcel2603/spikeball-league/internal/handler/league"
	staticfiles "github.com/Marcel2603/spikeball-league/internal/handler/static-files"
	"github.com/Marcel2603/spikeball-league/internal/handler/version"
	custommw "github.com/Marcel2603/spikeball-league/internal/middleware"
	"github.com/go-chi/chi/v5"

	"github.com/go-chi/cors"
	"github.com/go-chi/metrics"

	"github.com/go-chi/chi/v5/middleware"
)

//go:generate go tool github.com/a-h/templ/cmd/templ generate

//go:embed static
var staticFiles embed.FS

var (
	BuildVersion = "dev"
	BuildCommit  = "none"
)

func main() {
	configuration := config.Configuration

	var level slog.Level
	if err := level.UnmarshalText([]byte(configuration.Log.Level)); err != nil {
		level = slog.LevelInfo
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}))
	slog.SetDefault(logger)

	logger.Info("Starting SpikeBall League service",
		slog.String("version", BuildVersion),
		slog.String("commit", BuildCommit),
	)

	staticfiles.NewHandler(staticFiles)

	dbConn, err := db.InitDB(configuration.Database.Path)
	if err != nil {
		logger.Error("Failed to initialize database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer dbConn.Close()
	queries := db.New(dbConn)

	r, err := setupApp(configuration, logger, queries)
	if err != nil {
		logger.Error("Error starting server", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("Listening on :" + configuration.Server.Port)
	err = http.ListenAndServe(":"+configuration.Server.Port, r)
	if err != nil {
		logger.Error("Error starting server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func setupApp(configuration config.Config, logger *slog.Logger, queries *db.Queries) (*chi.Mux, error) {
	r := setupServerRouter(configuration, logger)

	domain := configuration.Server.Domain
	if domain == "" {
		domain = "https://" + configuration.Server.Host + ":" + configuration.Server.Port
	}
	leagueHandler := league.NewHandler(queries, domain)
	leagueHandler.RegisterRoutes(r)
	adminHandler := admin.NewHandler(queries, domain)
	adminHandler.RegisterRoutes(r)
	r.Get("/health", health.LivenessHandler)
	r.Get("/health/live", health.LivenessHandler)
	r.Get("/health/ready", health.ReadinessHandler())
	r.Handle("/metrics", metrics.Handler())
	r.Get("/version", version.Handler(BuildVersion, BuildCommit))

	r.Get("/favicon.ico", staticfiles.HandleFavicon)
	r.Get("/static/*", staticfiles.Handler)

	return r, nil
}

func setupServerRouter(configuration config.Config, logger *slog.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{configuration.Server.Host},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		ExposedHeaders:   []string{},
		AllowCredentials: false,
		MaxAge:           600,
	}))
	r.Use(metrics.Collector(metrics.CollectorOpts{
		Host:  false,
		Proto: true,
		Skip: func(r *http.Request) bool {
			if strings.HasPrefix(r.URL.Path, "/health") || r.URL.Path == "/metrics" {
				return true
			}
			return r.Method == http.MethodOptions
		},
	}))
	r.Use(middleware.RequestID)
	r.Use(custommw.SlogLogger(logger))
	r.Use(custommw.ThemeMiddleware)
	r.Use(middleware.Recoverer)
	return r
}
