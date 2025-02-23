package users

import "github.com/golang-jwt/jwt"

// Claims personalizados que incluirán la información del usuario
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
