package handler

import (
	"expenses-api/internal/api/dto"
	"expenses-api/internal/application/usecase"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ConfigHandler handles configuration-related HTTP requests
type ConfigHandler struct {
	salaryUseCase             *usecase.SalaryUseCase
	pocketUseCase             *usecase.PocketUseCase
	dailyExpenseConfigUseCase *usecase.DailyExpenseConfigUseCase
}

// NewConfigHandler creates a new config handler instance
func NewConfigHandler(
	salaryUseCase *usecase.SalaryUseCase,
	pocketUseCase *usecase.PocketUseCase,
	dailyExpenseConfigUseCase *usecase.DailyExpenseConfigUseCase,
) *ConfigHandler {
	return &ConfigHandler{
		salaryUseCase:             salaryUseCase,
		pocketUseCase:             pocketUseCase,
		dailyExpenseConfigUseCase: dailyExpenseConfigUseCase,
	}
}

// GetIncome obtiene la configuración de ingresos para un mes específico
// GET /api/config/income/{month}
// Implementa herencia automática del mes anterior si no existe configuración
func (h *ConfigHandler) GetIncome(c *gin.Context) {
	monthParam := c.Param("month")

	// Validate month format
	_, err := time.Parse("2006-01", monthParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid month format. Use YYYY-MM",
		})
		return
	}

	// Get salary with inheritance
	salary, err := h.salaryUseCase.GetByMonthWithInheritance(monthParam)
	if err != nil {
		// Si no hay configuración ni herencia, retornar valores por defecto
		response := dto.SalaryDTO{
			MonthlyAmount: 0,
		}
		c.JSON(http.StatusOK, response)
		return
	}

	response := dto.SalaryDTO{
		MonthlyAmount: salary.MonthlyAmount,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateIncome actualiza la configuración de ingresos para un mes específico
// PUT /api/config/income/{month}
func (h *ConfigHandler) UpdateIncome(c *gin.Context) {
	monthParam := c.Param("month")

	// Validate month format
	_, err := time.Parse("2006-01", monthParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid month format. Use YYYY-MM",
		})
		return
	}

	var salaryDTO dto.SalaryDTO
	if err := c.ShouldBindJSON(&salaryDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Update salary using use case for specified month
	err = h.salaryUseCase.UpdateSalary(salaryDTO.MonthlyAmount, monthParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error updating income configuration",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, salaryDTO)
}

// GetPockets obtiene todos los bolsillos
// GET /api/config/pockets
func (h *ConfigHandler) GetPockets(c *gin.Context) {
	pockets, err := h.pocketUseCase.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error getting pockets",
			"details": err.Error(),
		})
		return
	}

	// Convert to DTOs
	var pocketDTOs []dto.PocketDTO
	for _, pocket := range pockets {
		pocketDTOs = append(pocketDTOs, dto.PocketDTO{
			ID:          int(pocket.ID),
			Name:        pocket.Name,
			Description: pocket.Description,
		})
	}

	c.JSON(http.StatusOK, pocketDTOs)
}

// CreatePocket crea un nuevo bolsillo
// POST /api/config/pockets
func (h *ConfigHandler) CreatePocket(c *gin.Context) {
	var pocketDTO dto.PocketDTO
	if err := c.ShouldBindJSON(&pocketDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Create pocket using use case
	pocket, err := h.pocketUseCase.Create(pocketDTO.Name, pocketDTO.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error creating pocket",
			"details": err.Error(),
		})
		return
	}

	// Return created pocket as DTO
	responseDTO := dto.PocketDTO{
		ID:          int(pocket.ID),
		Name:        pocket.Name,
		Description: pocket.Description,
	}

	c.JSON(http.StatusCreated, responseDTO)
}

// UpdatePocket actualiza un bolsillo existente
// PUT /api/config/pockets/{id}
func (h *ConfigHandler) UpdatePocket(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid pocket ID",
		})
		return
	}

	var pocketDTO dto.PocketDTO
	if err := c.ShouldBindJSON(&pocketDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Update pocket using use case
	pocket, err := h.pocketUseCase.Update(uint(id), pocketDTO.Name, pocketDTO.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error updating pocket",
			"details": err.Error(),
		})
		return
	}

	// Return updated pocket as DTO
	responseDTO := dto.PocketDTO{
		ID:          int(pocket.ID),
		Name:        pocket.Name,
		Description: pocket.Description,
	}

	c.JSON(http.StatusOK, responseDTO)
}

// DeletePocket elimina un bolsillo
// DELETE /api/config/pockets/{id}
func (h *ConfigHandler) DeletePocket(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid pocket ID",
		})
		return
	}

	// Delete pocket using use case
	err = h.pocketUseCase.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error deleting pocket",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pocket deleted successfully",
		"id":      id,
	})
}

// GetDailyBudget obtiene la configuración de presupuesto diario para un mes específico
// GET /api/config/daily-budget/{month}
// Implementa herencia automática del mes anterior si no existe configuración
func (h *ConfigHandler) GetDailyBudget(c *gin.Context) {
	monthParam := c.Param("month")

	// Validate month format
	_, err := time.Parse("2006-01", monthParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid month format. Use YYYY-MM",
		})
		return
	}

	// Get daily expense config with inheritance
	config, err := h.dailyExpenseConfigUseCase.GetByMonthWithInheritance(monthParam)
	if err != nil {
		// Si no hay configuración ni herencia, retornar valores por defecto
		response := dto.DailyExpensesConfigDTO{
			MonthlyBudget: 0,
		}
		c.JSON(http.StatusOK, response)
		return
	}

	response := dto.DailyExpensesConfigDTO{
		MonthlyBudget: config.MonthlyBudget,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateDailyBudget actualiza la configuración de presupuesto diario para un mes específico
// PUT /api/config/daily-budget/{month}
func (h *ConfigHandler) UpdateDailyBudget(c *gin.Context) {
	monthParam := c.Param("month")

	// Validate month format
	_, err := time.Parse("2006-01", monthParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid month format. Use YYYY-MM",
		})
		return
	}

	var configDTO dto.DailyExpensesConfigDTO
	if err := c.ShouldBindJSON(&configDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Update daily budget using use case for specified month
	err = h.dailyExpenseConfigUseCase.UpdateBudget(configDTO.MonthlyBudget, monthParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error updating daily budget configuration",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, configDTO)
}
