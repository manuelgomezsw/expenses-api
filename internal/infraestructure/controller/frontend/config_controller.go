package frontend

import (
	"expenses-api/internal/api/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ConfigController maneja los endpoints de configuración para el frontend
type ConfigController struct {
	// TODO: Agregar use cases cuando estén conectados
}

// GetIncome obtiene la configuración de ingresos
// GET /api/config/income
func (cc *ConfigController) GetIncome(c *gin.Context) {
	// TODO: Usar salary use case real
	response := dto.SalaryDTO{
		MonthlyAmount: 0,
		Currency:      "COP",
	}

	c.JSON(http.StatusOK, response)
}

// UpdateIncome actualiza la configuración de ingresos
// PUT /api/config/income
func (cc *ConfigController) UpdateIncome(c *gin.Context) {
	var salaryDTO dto.SalaryDTO
	if err := c.ShouldBindJSON(&salaryDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// TODO: Convertir DTO a modelo de dominio y actualizar
	c.JSON(http.StatusOK, salaryDTO)
}

// GetPockets obtiene todos los bolsillos
// GET /api/config/pockets
func (cc *ConfigController) GetPockets(c *gin.Context) {
	// TODO: Implementar lógica real
	pockets := []dto.PocketDTO{
		{ID: 1, Name: "Comida", Budget: 500000, Color: "#FF6B6B"},
		{ID: 2, Name: "Transporte", Budget: 200000, Color: "#4ECDC4"},
		{ID: 3, Name: "Entretenimiento", Budget: 300000, Color: "#45B7D1"},
	}

	c.JSON(http.StatusOK, pockets)
}

// CreatePocket crea un nuevo bolsillo
// POST /api/config/pockets
func (cc *ConfigController) CreatePocket(c *gin.Context) {
	var pocketDTO dto.PocketDTO
	if err := c.ShouldBindJSON(&pocketDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// TODO: Implementar lógica real de creación
	pocketDTO.ID = 999 // ID temporal
	c.JSON(http.StatusCreated, pocketDTO)
}

// UpdatePocket actualiza un bolsillo existente
// PUT /api/config/pockets/{id}
func (cc *ConfigController) UpdatePocket(c *gin.Context) {
	id := c.Param("id")

	var pocketDTO dto.PocketDTO
	if err := c.ShouldBindJSON(&pocketDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// TODO: Implementar lógica real de actualización
	c.JSON(http.StatusOK, gin.H{
		"message": "Pocket updated successfully",
		"id":      id,
	})
}

// DeletePocket elimina un bolsillo
// DELETE /api/config/pockets/{id}
func (cc *ConfigController) DeletePocket(c *gin.Context) {
	id := c.Param("id")

	// TODO: Implementar lógica real de eliminación
	c.JSON(http.StatusOK, gin.H{
		"message": "Pocket deleted successfully",
		"id":      id,
	})
}

// NewConfigController crea una nueva instancia del controlador
func NewConfigController() *ConfigController {
	return &ConfigController{}
}
