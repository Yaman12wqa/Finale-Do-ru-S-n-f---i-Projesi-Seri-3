package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secscan/secscan/backend/internal/domain"
	"github.com/secscan/secscan/backend/internal/report"
	"github.com/secscan/secscan/backend/internal/service"
)

type ScanHandler struct {
	service *service.ScanService
}

func NewScanHandler(service *service.ScanService) *ScanHandler {
	return &ScanHandler{service: service}
}

func (h *ScanHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ScanHandler) StartScan(c *gin.Context) {
	var request domain.ScanRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	response, err := h.service.StartScan(c.Request.Context(), request)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusAccepted, response)
}

func (h *ScanHandler) GetScan(c *gin.Context) {
	scan, ok := h.service.GetScan(c.Param("id"))
	if !ok {
		writeError(c, http.StatusNotFound, "scan not found")
		return
	}
	c.JSON(http.StatusOK, scan)
}

func (h *ScanHandler) StreamScan(c *gin.Context) {
	scanID := c.Param("id")
	scan, ok := h.service.GetScan(scanID)
	if !ok {
		writeError(c, http.StatusNotFound, "scan not found")
		return
	}

	writer := c.Writer
	flusher, ok := writer.(http.Flusher)
	if !ok {
		writeError(c, http.StatusInternalServerError, "streaming is not supported")
		return
	}

	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")
	writer.Header().Set("X-Accel-Buffering", "no")

	if err := writeSSE(writer, flusher, "snapshot", scan); err != nil {
		return
	}

	if scan.Status == domain.ScanStatusCompleted || scan.Status == domain.ScanStatusFailed {
		return
	}

	events, unsubscribe := h.service.Subscribe(scanID)
	defer unsubscribe()

	for {
		select {
		case event, open := <-events:
			if !open {
				return
			}
			if err := writeSSE(writer, flusher, event.Type, event); err != nil {
				return
			}
			if event.Type == domain.EventScanCompleted || event.Type == domain.EventScanFailed {
				return
			}
		case <-c.Request.Context().Done():
			return
		}
	}
}

func (h *ScanHandler) DownloadReport(c *gin.Context) {
	scan, ok := h.service.GetScan(c.Param("id"))
	if !ok {
		writeError(c, http.StatusNotFound, "scan not found")
		return
	}
	if scan.Status != domain.ScanStatusCompleted {
		writeError(c, http.StatusConflict, "report is available after the scan completes")
		return
	}

	pdf, err := report.GeneratePDF(scan)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not generate report")
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"secscan-%s.pdf\"", scan.ID))
	c.Data(http.StatusOK, "application/pdf", pdf)
}

func writeSSE(writer gin.ResponseWriter, flusher http.Flusher, event string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(writer, "event: %s\ndata: %s\n\n", event, data); err != nil {
		return err
	}
	flusher.Flush()
	return nil
}

func writeError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"message": message,
			"status":  status,
		},
	})
}
