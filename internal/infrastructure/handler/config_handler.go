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
	salaryUseCase *usecase.SalaryUseCase
	pocketUseCase *usecase.PocketUseCase
}

// NewConfigHandler creates a new config handler instance
func NewConfigHandler(
	salaryUseCase *usecase.SalaryUseCase,
	pocketUseCase *usecase.PocketUseCase,
) *ConfigHandler {
	return &ConfigHandler{
		salaryUseCase: salaryUseCase,
		pocketUseCase: pocketUseCase,
	}
}

// GetIncome obtiene la configuración de ingresos para un mes específico
// GET /api/config/income/{month}
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

	// Get salary for specified month
	salary, err := h.salaryUseCase.GetByMonth(monthParam)
	if err != nil {
		// If no salary found, return default values
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
	pocket, err := h.pocketUseCase.Create(pocketDTO.Name, "")
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
		Budget:      pocketDTO.Budget,
		Description: pocketDTO.Description,
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
	pocket, err := h.pocketUseCase.Update(uint(id), pocketDTO.Name, "")
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
		Budget:      pocketDTO.Budget,
		Description: pocketDTO.Description,
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
