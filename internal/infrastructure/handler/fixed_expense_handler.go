package handler

import (
	"expenses-api/internal/api/dto"
	"expenses-api/internal/application/usecase"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// FixedExpenseHandler handles fixed expense-related HTTP requests
type FixedExpenseHandler struct {
	fixedExpenseUseCase *usecase.FixedExpenseUseCase
}

// NewFixedExpenseHandler creates a new fixed expense handler instance
func NewFixedExpenseHandler(fixedExpenseUseCase *usecase.FixedExpenseUseCase) *FixedExpenseHandler {
	return &FixedExpenseHandler{
		fixedExpenseUseCase: fixedExpenseUseCase,
	}
}

// GetByMonth obtiene los gastos fijos de un mes específico
// GET /api/fixed-expenses/{month}
func (h *FixedExpenseHandler) GetByMonth(c *gin.Context) {
	monthParam := c.Param("month")

	// Validate month format
	_, err := time.Parse("2006-01", monthParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid month format. Use YYYY-MM",
		})
		return
	}

	// Get fixed expenses using use case
	expenses, err := h.fixedExpenseUseCase.GetByMonth(monthParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error getting fixed expenses",
			"details": err.Error(),
		})
		return
	}

	// Convert to DTOs
	var expenseDTOs []dto.FixedExpenseDTO
	for _, expense := range expenses {
		var paidDate *time.Time
		if expense.PaidDate != nil {
			if parsedDate, err := time.Parse("2006-01-02", *expense.PaidDate); err == nil {
				paidDate = &parsedDate
			}
		}

		expenseDTOs = append(expenseDTOs, dto.FixedExpenseDTO{
			ID:        int(expense.ID),
			Name:      expense.ConceptName,
			Amount:    expense.Amount,
			DueDate:   expense.PaymentDay,
			IsPaid:    expense.IsPaid,
			PaidDate:  paidDate,
			CreatedAt: &expense.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, expenseDTOs)
}

// UpdateStatus actualiza el estado de pago de un gasto fijo
// PUT /api/fixed-expenses/{id}/status
func (h *FixedExpenseHandler) UpdateStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
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

	// Update payment status using use case
	err = h.fixedExpenseUseCase.UpdatePaymentStatus(uint(id), statusUpdate.IsPaid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error updating payment status",
			"details": err.Error(),
		})
		return
	}

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
