package handler

import (
	"expenses-api/internal/api/dto"
	"expenses-api/internal/application/usecase"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// DailyExpenseHandler handles daily expense-related HTTP requests
type DailyExpenseHandler struct {
	dailyExpenseUseCase *usecase.DailyExpenseUseCase
}

// NewDailyExpenseHandler creates a new daily expense handler instance
func NewDailyExpenseHandler(dailyExpenseUseCase *usecase.DailyExpenseUseCase) *DailyExpenseHandler {
	return &DailyExpenseHandler{
		dailyExpenseUseCase: dailyExpenseUseCase,
	}
}

// GetByMonth obtiene los gastos diarios de un mes específico
// GET /api/daily-expenses/{month}
func (h *DailyExpenseHandler) GetByMonth(c *gin.Context) {
	monthParam := c.Param("month")

	// Validate month format
	_, err := time.Parse("2006-01", monthParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid month format. Use YYYY-MM",
		})
		return
	}

	// Get daily expenses using use case
	expenses, err := h.dailyExpenseUseCase.GetByMonth(monthParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error getting daily expenses",
			"details": err.Error(),
		})
		return
	}

	// Convert to DTOs
	var expenseDTOs []dto.DailyExpenseDTO
	for _, expense := range expenses {
		expenseDTOs = append(expenseDTOs, dto.DailyExpenseDTO{
			ID:          int(expense.ID),
			Amount:      expense.Amount,
			Description: expense.Description,
			Date:        expense.Date,
			PocketID:    1, // TODO: Add pocket relationship if needed
		})
	}

	c.JSON(http.StatusOK, expenseDTOs)
}

// Create crea un nuevo gasto diario
// POST /api/daily-expenses
func (h *DailyExpenseHandler) Create(c *gin.Context) {
	var expenseDTO dto.DailyExpenseDTO
	if err := c.ShouldBindJSON(&expenseDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Create daily expense using use case
	expense, err := h.dailyExpenseUseCase.Create(
		expenseDTO.Description,
		expenseDTO.Amount,
		expenseDTO.Date,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error creating daily expense",
			"details": err.Error(),
		})
		return
	}

	// Return created expense as DTO
	responseDTO := dto.DailyExpenseDTO{
		ID:          int(expense.ID),
		Amount:      expense.Amount,
		Description: expense.Description,
		Date:        expense.Date,
		PocketID:    expenseDTO.PocketID,
	}

	c.JSON(http.StatusCreated, responseDTO)
}

// Update actualiza un gasto diario existente
// PUT /api/daily-expenses/{id}
func (h *DailyExpenseHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid expense ID",
		})
		return
	}

	var expenseDTO dto.DailyExpenseDTO
	if err := c.ShouldBindJSON(&expenseDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Update daily expense using use case
	expense, err := h.dailyExpenseUseCase.Update(
		uint(id),
		expenseDTO.Description,
		expenseDTO.Amount,
		expenseDTO.Date,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error updating daily expense",
			"details": err.Error(),
		})
		return
	}

	// Return updated expense as DTO
	responseDTO := dto.DailyExpenseDTO{
		ID:          int(expense.ID),
		Amount:      expense.Amount,
		Description: expense.Description,
		Date:        expense.Date,
		PocketID:    expenseDTO.PocketID,
	}

	c.JSON(http.StatusOK, responseDTO)
}

// Delete elimina un gasto diario
// DELETE /api/daily-expenses/{id}
func (h *DailyExpenseHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid expense ID",
		})
		return
	}

	// Delete daily expense using use case
	err = h.dailyExpenseUseCase.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error deleting daily expense",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Expense deleted successfully",
		"id":      id,
	})
}
