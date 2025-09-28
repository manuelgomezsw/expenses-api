package frontend

import (
	"expenses-api/internal/api/dto"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// FixedExpensesController maneja los endpoints de gastos fijos para el frontend
type FixedExpensesController struct{}

// GetByMonth obtiene los gastos fijos de un mes específico
// GET /api/fixed-expenses/{month}
func (fec *FixedExpensesController) GetByMonth(c *gin.Context) {
	monthParam := c.Param("month")

	// Validar formato de mes (YYYY-MM)
	_, err := time.Parse("2006-01", monthParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid month format. Use YYYY-MM",
		})
		return
	}

	// TODO: Implementar lógica real basada en concepts
	// Por ahora retornamos datos de ejemplo
	now := time.Now()
	fixedExpenses := []dto.FixedExpenseDTO{
		{
			ID:        1,
			Name:      "Arriendo",
			Amount:    1200000,
			DueDate:   5,
			IsPaid:    true,
			PaidDate:  &now,
			CreatedAt: &now,
		},
		{
			ID:        2,
			Name:      "Internet",
			Amount:    80000,
			DueDate:   15,
			IsPaid:    false,
			PaidDate:  nil,
			CreatedAt: &now,
		},
		{
			ID:        3,
			Name:      "Servicios",
			Amount:    200000,
			DueDate:   20,
			IsPaid:    false,
			PaidDate:  nil,
			CreatedAt: &now,
		},
	}

	c.JSON(http.StatusOK, fixedExpenses)
}

// UpdateStatus actualiza el estado de pago de un gasto fijo
// PUT /api/fixed-expenses/{id}/status
func (fec *FixedExpensesController) UpdateStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid expense ID",
		})
		return
	}

	var statusUpdate struct {
		IsPaid bool `json:"is_paid" binding:"required"`
	}

	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// TODO: Implementar lógica real de actualización de estado
	// Esto debería actualizar el campo 'payed' en la tabla concepts

	response := gin.H{
		"message": "Status updated successfully",
		"id":      id,
		"is_paid": statusUpdate.IsPaid,
	}

	if statusUpdate.IsPaid {
		now := time.Now()
		response["paid_date"] = now
	}

	c.JSON(http.StatusOK, response)
}

// NewFixedExpensesController crea una nueva instancia del controlador
func NewFixedExpensesController() *FixedExpensesController {
	return &FixedExpensesController{}
}
