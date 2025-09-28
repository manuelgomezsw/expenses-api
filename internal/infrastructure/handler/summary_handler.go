package handler

import (
	"expenses-api/internal/application/usecase"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SummaryHandler handles summary-related HTTP requests
type SummaryHandler struct {
	summaryUseCase *usecase.SummaryUseCase
}

// NewSummaryHandler creates a new summary handler instance
func NewSummaryHandler(summaryUseCase *usecase.SummaryUseCase) *SummaryHandler {
	return &SummaryHandler{
		summaryUseCase: summaryUseCase,
	}
}

// GetMonthlySummary obtiene el resumen financiero mensual
// GET /api/summary/{month}
func (h *SummaryHandler) GetMonthlySummary(c *gin.Context) {
	monthParam := c.Param("month")

	// Validate month format
	_, err := time.Parse("2006-01", monthParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid month format. Use YYYY-MM",
		})
		return
	}

	// Get monthly summary using use case
	summary, err := h.summaryUseCase.GetMonthlySummary(monthParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error calculating monthly summary",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, summary)
}
