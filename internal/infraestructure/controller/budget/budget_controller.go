package budget

import (
	"expenses-api/internal/domain/budgets/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Calculate(c *gin.Context) {
	cycleID, err := strconv.ParseInt(c.Param("cycle_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error serializing cycle_id",
			"error":   err.Error(),
		})
		return
	}

	budget, err := service.Calculate(int(cycleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error calculating budget",
			"error":   err.Error(),
		})
		return
	}

	if budget.CycleID == 0 {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, budget)
}
