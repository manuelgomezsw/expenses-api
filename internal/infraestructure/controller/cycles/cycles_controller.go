package cycles

import (
	"expenses-api/internal/domain/cycles"
	"expenses-api/internal/domain/cycles/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetAll(c *gin.Context) {
	allCycles, err := service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error getting all cycles",
			"error":   err.Error(),
		})
		return
	}

	if allCycles == nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, allCycles)
}

func GetActive(c *gin.Context) {
	allCycles, err := service.GetActive()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error getting all active cycles",
			"error":   err.Error(),
		})
		return
	}

	if allCycles == nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, allCycles)
}

func GetByID(c *gin.Context) {
	cycleID, err := strconv.ParseInt(c.Param("cycle_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error serializing cycle_id",
			"error":   err.Error(),
		})
		return
	}

	cycle, err := service.GetByID(int(cycleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error getting cycle",
			"error":   err.Error(),
		})
		return
	}

	if cycle.CycleID == 0 {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, cycle)
}

func Create(c *gin.Context) {
	var cycle cycles.Cycle
	if err := c.ShouldBindJSON(&cycle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error serializing body",
			"error":   err.Error(),
		})
		return
	}

	if err := service.Create(&cycle); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error creating cycle",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, cycle)
}

func Update(c *gin.Context) {
	cycleID, err := strconv.ParseInt(c.Param("cycle_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error serializing cycle_id",
			"error":   err.Error(),
		})
		return
	}

	var cycle cycles.Cycle
	if err := c.ShouldBindJSON(&cycle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error serializing body",
			"error":   err.Error(),
		})
		return
	}

	if err := service.Update(int(cycleID), &cycle); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error updating cycle",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, cycle)
}

func Delete(c *gin.Context) {
	cycleID, err := strconv.ParseInt(c.Param("cycle_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error serializing cycle_id",
			"error":   err.Error(),
		})
		return
	}

	if err := service.Delete(int(cycleID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error deleting cycle",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, nil)
}
