package frontend

import (
	"expenses-api/internal/api/dto"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// DailyExpensesController maneja los endpoints de gastos diarios para el frontend
type DailyExpensesController struct{}

// GetByMonth obtiene los gastos diarios de un mes específico
// GET /api/daily-expenses/{month}
func (dec *DailyExpensesController) GetByMonth(c *gin.Context) {
	monthParam := c.Param("month")

	// Validar formato de mes (YYYY-MM)
	_, err := time.Parse("2006-01", monthParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid month format. Use YYYY-MM",
		})
		return
	}

	// TODO: Implementar lógica real basada en expenses
	// Por ahora retornamos datos de ejemplo
	dailyExpenses := []dto.DailyExpenseDTO{
		{
			ID:          1,
			Amount:      25000,
			Description: "Almuerzo",
			Date:        "2024-01-15",
			PocketID:    1, // Comida
		},
		{
			ID:          2,
			Amount:      8000,
			Description: "Transporte público",
			Date:        "2024-01-15",
			PocketID:    2, // Transporte
		},
		{
			ID:          3,
			Amount:      45000,
			Description: "Cine",
			Date:        "2024-01-14",
			PocketID:    3, // Entretenimiento
		},
	}

	c.JSON(http.StatusOK, dailyExpenses)
}

// Create crea un nuevo gasto diario
// POST /api/daily-expenses
func (dec *DailyExpensesController) Create(c *gin.Context) {
	var expenseDTO dto.DailyExpenseDTO
	if err := c.ShouldBindJSON(&expenseDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validar fecha
	_, err := time.Parse("2006-01-02", expenseDTO.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid date format. Use YYYY-MM-DD",
		})
		return
	}

	// TODO: Implementar lógica real de creación
	expenseDTO.ID = 999 // ID temporal
	c.JSON(http.StatusCreated, expenseDTO)
}

// Update actualiza un gasto diario existente
// PUT /api/daily-expenses/{id}
func (dec *DailyExpensesController) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
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

	// Validar fecha
	_, err = time.Parse("2006-01-02", expenseDTO.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid date format. Use YYYY-MM-DD",
		})
		return
	}

	// TODO: Implementar lógica real de actualización
	expenseDTO.ID = id
	c.JSON(http.StatusOK, expenseDTO)
}

// Delete elimina un gasto diario
// DELETE /api/daily-expenses/{id}
func (dec *DailyExpensesController) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid expense ID",
		})
		return
	}

	// TODO: Implementar lógica real de eliminación
	c.JSON(http.StatusOK, gin.H{
		"message": "Expense deleted successfully",
		"id":      id,
	})
}

// NewDailyExpensesController crea una nueva instancia del controlador
func NewDailyExpensesController() *DailyExpensesController {
	return &DailyExpensesController{}
}
