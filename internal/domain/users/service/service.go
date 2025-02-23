package service

import (
	"errors"
	"expenses-api/internal/domain/users"
	"expenses-api/internal/domain/users/repository"
	"expenses-api/internal/util/customdate"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	jwtKey = []byte("mi_clave_secreta_super_segura")
)

func GetByID(userID int) (users.UserOutput, error) {
	return repository.GetByID(userID)
}

func GetByUsername(username string) (users.UserOutput, error) {
	return repository.GetByUsername(username)
}

func Create(user users.UserInput) (users.UserOutput, error) {
	passwordHashed, err := hashPassword(user.Password)
	if err != nil {
		return users.UserOutput{}, err
	}
	user.Password = passwordHashed

	createdAtNow, err := customdate.ParseAndFormatDateMySql(time.Now().String())
	if err != nil {
		return users.UserOutput{}, errors.New("fecha inicial inválida: " + err.Error())
	}
	user.CreatedAt = createdAtNow

	return repository.Create(user)
}

func Update(userID int, user users.UserInput) (users.UserOutput, error) {
	return repository.Update(userID, user)
}

func ChangePassword(username, currentPassword, newPassword string) (bool, error) {
	return true, nil
}

func Login(credentials users.Credentials) (users.TokenPair, error) {
	storedPassword, err := repository.GetPasswordByUsername(credentials.Username)
	if err != nil {
		return users.TokenPair{}, err
	}

	if err = checkPassword(storedPassword, credentials.Password); err != nil {
		return users.TokenPair{}, err
	}

	// Definir tiempos de expiración: 5 minutos para el access token y 24 horas para el refresh token
	accessExpirationTime := time.Now().Add(2 * time.Minute)
	refreshExpirationTime := time.Now().Add(3 * time.Hour)

	// Crear claims para el access token
	accessClaims := &users.Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessExpirationTime.Unix(),
		},
	}
	// Crear claims para el refresh token
	refreshClaims := &users.Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshExpirationTime.Unix(),
		},
	}

	// Crear el access token usando HS256
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		return users.TokenPair{}, err
	}

	// Crear el refresh token usando HS256
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		return users.TokenPair{}, err
	}

	return users.TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func RefreshToken(oldRefreshToken string) (string, error) {
	// Parsear el refresh token usando los mismos claims que en login
	token, err := jwt.ParseWithClaims(oldRefreshToken, &users.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*users.Claims)
	if !ok || !token.Valid {
		return "", errors.New("refresh token no válido")
	}

	// Generar un nuevo access token con expiración de 5 minutos
	newExpirationTime := time.Now().Add(5 * time.Minute)
	newClaims := &users.Claims{
		Username: claims.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: newExpirationTime.Unix(),
		},
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	newTokenString, err := newToken.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return newTokenString, nil
}

// HashPassword toma una contraseña en texto plano y la convierte en un hash bcrypt.
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword compara la contraseña ingresada con el hash almacenado.
func checkPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
