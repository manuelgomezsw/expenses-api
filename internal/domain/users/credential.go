package users

// Credentials Estructura para recibir las credenciales del login
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
