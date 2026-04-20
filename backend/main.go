package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/secscan/secscan/backend/internal/api/routes"
	"github.com/secscan/secscan/backend/internal/config"
	"github.com/secscan/secscan/backend/internal/scanner/cve"
	"github.com/secscan/secscan/backend/internal/scanner/fuzz"
	"github.com/secscan/secscan/backend/internal/scanner/headers"
	"github.com/secscan/secscan/backend/internal/scanner/ports"
	"github.com/secscan/secscan/backend/internal/scanner/registry"
	"github.com/secscan/secscan/backend/internal/scanner/sqli"
	tlsscanner "github.com/secscan/secscan/backend/internal/scanner/tls"
	"github.com/secscan/secscan/backend/internal/scanner/xss"
	"github.com/secscan/secscan/backend/internal/service"
	"github.com/secscan/secscan/backend/internal/sse"
	"github.com/secscan/secscan/backend/internal/storage"
)

func main() {
	cfg := config.Load()
	gin.SetMode(cfg.GinMode)

	store := storage.NewMemoryStore()
	broker := sse.NewBroker()
	reg := registry.New()
	reg.Register(ports.New())
	reg.Register(headers.New())
	reg.Register(tlsscanner.New())
	reg.Register(fuzz.New())
	reg.Register(xss.New())
	reg.Register(sqli.New())
	reg.Register(cve.New())

	scanService := service.NewScanService(store, broker, reg, cfg.ScanTimeout)

	router := gin.New()
	router.Use(gin.Recovery())
	routes.Register(router, scanService, cfg)

	addr := ":" + cfg.Port
	log.Printf("SecScan backend listening on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
