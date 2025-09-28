package frontend

import (
	"expenses-api/internal/api/dto"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SummaryController maneja los endpoints de resumen para el frontend
type SummaryController struct {
	// TODO: Agregar use cases cuando estén conectados
}

// GetMonthlySummary obtiene el resumen financiero mensual
// GET /api/summary/{month}
func (sc *SummaryController) GetMonthlySummary(c *gin.Context) {
	monthParam := c.Param("month")

	// Validar formato de mes (YYYY-MM)
	_, err := time.Parse("2006-01", monthParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid month format. Use YYYY-MM",
		})
		return
	}

	// TODO: Implementar lógica real de cálculo de resumen
	// Por ahora retornamos datos de ejemplo
	summary := dto.MonthlySummaryDTO{
		Month:              monthParam,
		TotalIncome:        5000000, // 5M COP
		TotalFixedExpenses: 2500000, // 2.5M COP
		TotalDailyExpenses: 800000,  // 800K COP
		RemainingBudget:    1700000, // 1.7M COP
		FixedExpensesPaid:  8,
		FixedExpensesTotal: 12,
		DailyBudgetUsed:    800000,
		DailyBudgetTotal:   1500000,
	}

	c.JSON(http.StatusOK, summary)
}

// NewSummaryController crea una nueva instancia del controlador
func NewSummaryController() *SummaryController {
	return &SummaryController{}
}
