package pockets

import (
	"expenses-api/internal/domain/pockets"
	"expenses-api/internal/domain/pockets/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Get(c *gin.Context) {
	allPockets, err := service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error getting all pockets",
			"error":   err.Error(),
		})
		return
	}

	if allPockets == nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, allPockets)
}

func GetActives(c *gin.Context) {
	allPockets, err := service.GetActives()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error getting all pockets",
			"error":   err.Error(),
		})
		return
	}

	if allPockets == nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, allPockets)
}

func GetByID(c *gin.Context) {
	pocketID, err := strconv.ParseInt(c.Param("pocket_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error serializing pocket_id",
			"error":   err.Error(),
		})
		return
	}

	pocket, err := service.GetByID(pocketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error getting pocket",
			"error":   err.Error(),
		})
		return
	}

	if pocket.PocketID == 0 {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, pocket)
}

func Create(c *gin.Context) {
	var newPocket pockets.Pocket
	if err := c.ShouldBindJSON(&newPocket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error serializing body",
			"error":   err.Error(),
		})
		return
	}

	if err := service.Create(&newPocket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error creating pocket",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, newPocket)
}

func Update(c *gin.Context) {
	var currentPocket pockets.Pocket
	if err := c.ShouldBindJSON(&currentPocket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error serializing body",
			"error":   err.Error(),
		})
		return
	}

	pocketID, err := strconv.ParseInt(c.Param("pocket_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error serializing pocket_id",
			"error":   err.Error(),
		})
		return
	}

	if err := service.Update(pocketID, &currentPocket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error updating pocket",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, currentPocket)
}

func Delete(c *gin.Context) {
	pocketID, err := strconv.ParseInt(c.Param("pocket_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error serializing pocket_id",
			"error":   err.Error(),
		})
		return
	}

	if err := service.Delete(pocketID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error deleting pocket",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, nil)
}
