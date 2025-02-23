package concepts

import (
	"expenses-api/internal/domain/concepts"
	"expenses-api/internal/domain/concepts/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetByID(c *gin.Context) {
	conceptID, err := strconv.ParseInt(c.Param("concept_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error serializing concept_id",
			"error":   err.Error(),
		})
		return
	}

	conceptsByPocket, err := service.GetByID(int(conceptID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error getting concepts by id",
			"error":   err.Error(),
		})
		return
	}

	if len(conceptsByPocket) == 0 {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, conceptsByPocket)
}

func GetByPocketID(c *gin.Context) {
	pocketID, err := strconv.ParseInt(c.Param("pocket_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error serializing pocket_id",
			"error":   err.Error(),
		})
		return
	}

	conceptsByPocket, err := service.GetByPocketID(int(pocketID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error getting concepts by pocket",
			"error":   err.Error(),
		})
		return
	}

	if len(conceptsByPocket) == 0 {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, conceptsByPocket)
}

func Create(c *gin.Context) {
	var concept concepts.Concept
	if err := c.ShouldBindJSON(&concept); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error serializing body",
			"error":   err.Error(),
		})
		return
	}

	if err := service.Create(&concept); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error creating concept",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, concept)
}

func Update(c *gin.Context) {
	conceptID, err := strconv.ParseInt(c.Param("concept_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error serializing concept_id",
			"error":   err.Error(),
		})
		return
	}

	var concept concepts.Concept
	if err := c.ShouldBindJSON(&concept); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error serializing body",
			"error":   err.Error(),
		})
		return
	}

	if err := service.Update(int(conceptID), &concept); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error updating concept",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, concept)
}

func Delete(c *gin.Context) {
	conceptID, err := strconv.ParseInt(c.Param("concept_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error serializing concept_id",
			"error":   err.Error(),
		})
		return
	}

	if err := service.Delete(int(conceptID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error deleting concept",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, nil)
}
