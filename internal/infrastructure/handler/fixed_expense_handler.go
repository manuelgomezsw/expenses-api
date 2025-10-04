package handler

import (
	"expenses-api/internal/api/dto"
	"expenses-api/internal/application/usecase"
	"expenses-api/internal/domain/fixed_expense"
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
// Implementa herencia automática del mes anterior si no existen gastos
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

	// Get fixed expenses with inheritance
	expenses, err := h.fixedExpenseUseCase.GetByMonthWithInheritance(monthParam)
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
		// Get pocket name from the preloaded relationship
		pocketName := ""
		if expense.Pocket != nil {
			pocketName = expense.Pocket.Name
		}

		expenseDTOs = append(expenseDTOs, dto.FixedExpenseDTO{
			ID:          int(expense.ID), // Será 0 si es heredado
			PocketName:  pocketName,
			ConceptName: expense.ConceptName,
			Amount:      expense.Amount,
			PaymentDay:  expense.PaymentDay,
			Month:       expense.Month,    // Siempre el mes solicitado
			IsPaid:      expense.IsPaid,   // Será false si es heredado
			PaidDate:    expense.PaidDate, // Será nil si es heredado
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
		IsPaid *bool `json:"is_paid" binding:"required"`
	}

	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate that is_paid was provided
	if statusUpdate.IsPaid == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "is_paid field is required",
		})
		return
	}

	// Update payment status using use case
	err = h.fixedExpenseUseCase.UpdatePaymentStatus(uint(id), *statusUpdate.IsPaid)
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
		"is_paid": *statusUpdate.IsPaid,
	}

	if *statusUpdate.IsPaid {
		now := time.Now()
		response["paid_date"] = now
	}

	c.JSON(http.StatusOK, response)
}

// Create crea un nuevo gasto fijo
// POST /api/fixed-expenses
func (h *FixedExpenseHandler) Create(c *gin.Context) {
	var expenseDTO dto.FixedExpenseDTO

	if err := c.ShouldBindJSON(&expenseDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Convert DTO to domain model
	expense := &fixed_expense.FixedExpense{
		ConceptName: expenseDTO.ConceptName,
		Amount:      expenseDTO.Amount,
		PaymentDay:  expenseDTO.PaymentDay,
		PocketID:    uint(expenseDTO.PocketID),
		Month:       time.Now().Format("2006-01"), // Siempre usar mes actual
		IsPaid:      false,                        // Siempre false por defecto
	}

	// Create expense using use case
	err := h.fixedExpenseUseCase.Create(expense)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error creating fixed expense",
			"details": err.Error(),
		})
		return
	}

	// Get created expense with pocket information
	createdExpense, err := h.fixedExpenseUseCase.GetByID(expense.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error retrieving created expense",
			"details": err.Error(),
		})
		return
	}

	// Get pocket name from the preloaded relationship
	pocketName := ""
	if createdExpense.Pocket != nil {
		pocketName = createdExpense.Pocket.Name
	}

	// Convert back to DTO for response
	responseDTO := dto.FixedExpenseDTO{
		ID:          int(createdExpense.ID),
		PocketName:  pocketName,
		ConceptName: createdExpense.ConceptName,
		Amount:      createdExpense.Amount,
		PaymentDay:  createdExpense.PaymentDay,
		Month:       createdExpense.Month,
		IsPaid:      createdExpense.IsPaid,
		PaidDate:    createdExpense.PaidDate,
	}

	c.JSON(http.StatusCreated, responseDTO)
}

// Update actualiza un gasto fijo existente
// PUT /api/fixed-expenses/{id}
func (h *FixedExpenseHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid expense ID",
		})
		return
	}

	var expenseDTO dto.FixedExpenseDTO
	if err := c.ShouldBindJSON(&expenseDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Convert DTO to domain model
	updatedExpense := &fixed_expense.FixedExpense{
		ConceptName: expenseDTO.ConceptName,
		Amount:      expenseDTO.Amount,
		PaymentDay:  expenseDTO.PaymentDay,
		PocketID:    uint(expenseDTO.PocketID),
		Month:       time.Now().Format("2006-01"), // Siempre usar mes actual
	}

	// Update expense using use case
	err = h.fixedExpenseUseCase.Update(uint(id), updatedExpense)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "expense not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "expense ID is required" ||
			err.Error() == "expense data is required" ||
			err.Error() == "concept name is required" ||
			err.Error() == "amount must be greater than 0" ||
			err.Error() == "payment day must be between 1 and 31" ||
			err.Error() == "month is required" ||
			err.Error() == "pocket ID is required" ||
			err.Error() == "invalid month format, must be YYYY-MM" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error":   "Error updating fixed expense",
			"details": err.Error(),
		})
		return
	}

	// Get updated expense to return in response
	updatedExpenseFromDB, err := h.fixedExpenseUseCase.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error retrieving updated expense",
			"details": err.Error(),
		})
		return
	}

	// Get pocket name from the preloaded relationship
	pocketName := ""
	if updatedExpenseFromDB.Pocket != nil {
		pocketName = updatedExpenseFromDB.Pocket.Name
	}

	// Convert to DTO for response
	responseDTO := dto.FixedExpenseDTO{
		ID:          int(updatedExpenseFromDB.ID),
		PocketName:  pocketName,
		ConceptName: updatedExpenseFromDB.ConceptName,
		Amount:      updatedExpenseFromDB.Amount,
		PaymentDay:  updatedExpenseFromDB.PaymentDay,
		Month:       updatedExpenseFromDB.Month,
		IsPaid:      updatedExpenseFromDB.IsPaid,
		PaidDate:    updatedExpenseFromDB.PaidDate,
	}

	c.JSON(http.StatusOK, responseDTO)
}
