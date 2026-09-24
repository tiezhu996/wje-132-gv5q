package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"safetyplatform/internal/config"
	"safetyplatform/internal/handler"
	"safetyplatform/internal/model"
	"safetyplatform/internal/repository"
	"safetyplatform/internal/router"
	"safetyplatform/internal/service"
	"safetyplatform/internal/util"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()
	logger := util.NewLogger(slog.LevelInfo)

	db, err := gorm.Open(mysql.Open(cfg.DBDSN()), &gorm.Config{})
	if err != nil {
		logger.Error("connect database failed", "error", err.Error())
		os.Exit(1)
	}
	if err := db.AutoMigrate(
		&model.User{}, &model.SafetyIncident{}, &model.SafetyInspection{}, &model.InspectionItem{},
		&model.SafetyTraining{}, &model.WorkerCertification{}, &model.AuditLog{},
	); err != nil {
		logger.Error("auto migrate failed", "error", err.Error())
		os.Exit(1)
	}
	if err := service.NewSeedService(db, logger).Seed(); err != nil {
		logger.Error("seed failed", "error", err.Error())
		os.Exit(1)
	}

	userRepo := repository.NewUserRepository(db)
	incidentRepo := repository.NewSafetyIncidentRepository(db)
	inspectionRepo := repository.NewSafetyInspectionRepository(db)
	itemRepo := repository.NewInspectionItemRepository(db)
	trainingRepo := repository.NewSafetyTrainingRepository(db)
	certRepo := repository.NewWorkerCertificationRepository(db)

	userSvc := service.NewUserService(userRepo, logger)
	incidentSvc := service.NewSafetyIncidentService(incidentRepo, userRepo, logger)
	inspectionSvc := service.NewSafetyInspectionService(db, inspectionRepo, itemRepo, userRepo, logger)
	trainingSvc := service.NewSafetyTrainingService(trainingRepo, userRepo, logger)
	certSvc := service.NewWorkerCertificationService(certRepo, userRepo, logger)
	dashboardSvc := service.NewDashboardService(incidentSvc, inspectionSvc, trainingSvc, certSvc, logger)

	userHandler := handler.NewUserHandler(userSvc, logger)
	incidentHandler := handler.NewSafetyIncidentHandler(incidentSvc, logger)
	inspectionHandler := handler.NewSafetyInspectionHandler(inspectionSvc, logger)
	itemHandler := handler.NewInspectionItemHandler(inspectionSvc, logger)
	trainingHandler := handler.NewSafetyTrainingHandler(trainingSvc, logger)
	certHandler := handler.NewWorkerCertificationHandler(certSvc, logger)
	dashboardHandler := handler.NewDashboardHandler(dashboardSvc, logger)
	uploadHandler := handler.NewUploadHandler(cfg, logger)
	auditLogHandler := handler.NewAuditLogHandler(db, logger)

	r := router.New(cfg, db, logger, userHandler, incidentHandler, inspectionHandler, itemHandler,
		trainingHandler, certHandler, dashboardHandler, uploadHandler, auditLogHandler)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r.Setup(),
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("server starting", "port", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server run failed", "error", err.Error())
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("server shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server forced to shutdown", "error", err.Error())
	}
	logger.Info("server stopped")
}
