package expenses

import (
	"expenses-api/internal/domain/expenses"
	"expenses-api/internal/domain/expenses/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetByActiveCycles(c *gin.Context) {
	allExpenses, err := service.GetByActiveCycles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error getting expenses by active cycle",
			"error":   err.Error(),
		})
		return
	}

	if allExpenses == nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, allExpenses)
}

func GetByPocketID(c *gin.Context) {
	pocketID, err := strconv.ParseInt(c.Param("pocket_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error serializing pocket_id",
			"error":   err.Error(),
		})
		return
	}

	allExpenses, err := service.GetByPocketID(int(pocketID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error getting expenses from pocket",
			"error":   err.Error(),
		})
		return
	}

	if allExpenses == nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, allExpenses)
}

func GetByPaymentTypeID(c *gin.Context) {
	paymentTypeID, err := strconv.ParseInt(c.Param("payment_type_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error serializing payment_type_id",
			"error":   err.Error(),
		})
		return
	}

	allExpenses, err := service.GetByPaymentTypeID(int16(paymentTypeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error getting expenses from payment_type",
			"error":   err.Error(),
		})
		return
	}

	if allExpenses == nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, allExpenses)
}

func GetByID(c *gin.Context) {
	expenseID, err := strconv.ParseInt(c.Param("expense_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error serializing expense_id",
			"error":   err.Error(),
		})
		return
	}

	expense, err := service.GetByID(int(expenseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error getting expense",
			"error":   err.Error(),
		})
		return
	}

	if expense.PaymentTypeID == 0 {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, expense)
}

func Create(c *gin.Context) {
	var expense expenses.Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error serializing body",
			"error":   err.Error(),
		})
		return
	}

	if err := service.Create(&expense); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error creating payment type",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, expense)
}

func Update(c *gin.Context) {
	var expense expenses.Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error serializing body",
			"error":   err.Error(),
		})
		return
	}

	expenseID, err := strconv.ParseInt(c.Param("expense_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error serializing expense_id",
			"error":   err.Error(),
		})
		return
	}

	if err := service.Update(int(expenseID), &expense); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error updating expense",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, expense)
}

func Delete(c *gin.Context) {
	expenseID, err := strconv.ParseInt(c.Param("expense_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error serializing expense_id",
			"error":   err.Error(),
		})
		return
	}

	if err := service.Delete(int(expenseID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error deleting expense",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, nil)
}
