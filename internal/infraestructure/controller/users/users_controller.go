package users

import (
	"expenses-api/internal/domain/users"
	"expenses-api/internal/domain/users/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Create(c *gin.Context) {
	var user users.UserInput
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error serializing body",
			"error":   err.Error(),
		})
		return
	}

	userOutput, err := service.Create(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error creating user",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, userOutput)
}

func Login(c *gin.Context) {
	var credentials users.Credentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error serializing body",
			"error":   err.Error(),
		})
		return
	}

	token, err := service.Login(credentials)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "user unauthorized",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func Refresh(c *gin.Context) {
	var oldRefreshToken string
	if err := c.ShouldBindJSON(&oldRefreshToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error serializing body",
			"error":   err.Error(),
		})
		return
	}

	tokenRefreshed, err := service.RefreshToken(oldRefreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "refresh token invalid",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenRefreshed,
	})
}
