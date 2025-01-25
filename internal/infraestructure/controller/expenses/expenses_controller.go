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
