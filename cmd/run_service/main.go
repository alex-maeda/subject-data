// Command run_service runs the subject-data API server.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/sovraai/subject-data/internal/api"
	"github.com/sovraai/subject-data/internal/persistence"
	"github.com/sovraai/subject-data/internal/service"
	_ "github.com/sovraai/subject-data/internal/swaggerdocs" // swagger generated docs
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8380"
	}

	// Database setup
	dbDriver := os.Getenv("DB_DRIVER")
	if dbDriver == "" {
		dbDriver = "sqlite"
	}

	var db *sqlx.DB
	var err error
	switch dbDriver {
	case "sqlite":
		dbURL := os.Getenv("DATABASE_URL")
		if dbURL == "" {
			dbURL = "subject-data.db"
		}
		db, err = persistence.NewSQLiteDB(dbURL)
	case "postgres":
		db, err = persistence.NewPostgresDB(os.Getenv("DATABASE_URL"))
	default:
		logger.Error("unsupported DB_DRIVER", "driver", dbDriver)
		os.Exit(1)
	}
	if err != nil {
		logger.Error("database", "err", err)
		os.Exit(1)
	}
	logger.Info("database connected", "driver", dbDriver)

	// Run migrations
	migrateDriver := dbDriver
	if migrateDriver == "postgres" {
		migrateDriver = "pgx"
	}
	if err := persistence.RunMigrations(db, migrateDriver); err != nil {
		logger.Error("migrations", "err", err)
		os.Exit(1)
	}
	logger.Info("migrations applied")

	// Construct repos and services
	subjectRepo := persistence.NewSubjectRepo(db)
	recordsRepo := persistence.NewSubjectIndexedRepository(db, "records")
	ratingsRepo := persistence.NewSubjectIndexedRepository(db, "subject_ratings")
	featuresRepo := persistence.NewSubjectIndexedRepository(db, "subject_ce_features")

	subjectSvc := service.NewSubjectService(logger, subjectRepo, recordsRepo, ratingsRepo, featuresRepo)
	recordSvc := service.NewRecordService(logger, db)

	router := api.NewRouter(api.RouterDeps{
		Logger:         logger,
		AuthTokens:     strings.Split(os.Getenv("AUTH_TOKENS"), ","),
		CORSOrigins:    os.Getenv("CORS_ALLOW_ORIGINS"),
		SubjectService: subjectSvc,
		RecordService:  recordSvc,
	})

	srv := api.NewServer(router, port, logger)

	go func() {
		if err := srv.Start(); err != nil {
			logger.Error("server", "err", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("shutdown", "err", err)
	}
}
