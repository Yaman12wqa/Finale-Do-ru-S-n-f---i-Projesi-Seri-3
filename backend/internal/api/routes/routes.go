package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/secscan/secscan/backend/internal/api/handlers"
	"github.com/secscan/secscan/backend/internal/api/middleware"
	"github.com/secscan/secscan/backend/internal/config"
	"github.com/secscan/secscan/backend/internal/service"
)

func Register(router *gin.Engine, scanService *service.ScanService, cfg config.Config) {
	router.Use(middleware.CORS(cfg.FrontendOrigin))

	handler := handlers.NewScanHandler(scanService)

	router.GET("/health", handler.Health)

	api := router.Group("/api")
	api.POST("/scan", handler.StartScan)
	api.GET("/scan/:id", handler.GetScan)
	api.GET("/scan/:id/stream", handler.StreamScan)
	api.GET("/scan/:id/report.pdf", handler.DownloadReport)
}
