package paymentstype

import (
	"expenses-api/internal/domain/paymentstype"
	"expenses-api/internal/domain/paymentstype/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Get(c *gin.Context) {
	paymentsType, err := service.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error getting payments type",
			"error":   err.Error(),
		})
		return
	}

	if paymentsType == nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, paymentsType)
}

func GetByID(c *gin.Context) {
	paymentTypeID, err := strconv.ParseInt(c.Param("payment_type_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error serializing payment_type_id",
			"error":   err.Error(),
		})
		return
	}

	paymentType, err := service.GetByID(int16(paymentTypeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error getting payment type",
			"error":   err.Error(),
		})
		return
	}

	if paymentType.PaymentTypeID == 0 {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, paymentType)
}

func Create(c *gin.Context) {
	var newPaymentType paymentstype.PaymentType
	if err := c.ShouldBindJSON(&newPaymentType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error serializing body",
			"error":   err.Error(),
		})
		return
	}

	if err := service.Create(&newPaymentType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error creating payment type",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, newPaymentType)
}

func Update(c *gin.Context) {
	var paymentType paymentstype.PaymentType
	if err := c.ShouldBindJSON(&paymentType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error serializing body",
			"error":   err.Error(),
		})
		return
	}

	paymentTypeID, err := strconv.ParseInt(c.Param("payment_type_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error serializing payment_type_id",
			"error":   err.Error(),
		})
		return
	}

	if err := service.Update(int16(paymentTypeID), &paymentType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error updating payment type",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, paymentType)
}

func Delete(c *gin.Context) {
	paymentTypeId, err := strconv.ParseInt(c.Param("payment_type_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error serializing payment_type_id",
			"error":   err.Error(),
		})
		return
	}

	if err := service.Delete(int16(paymentTypeId)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error deleting payment type",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, nil)
}
