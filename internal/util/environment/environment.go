package environment

import "os"

const (
	Production = "production"
)

func IsProductionEnv() bool {
	appEnv := os.Getenv("APP_ENV")

	if appEnv == "" {
		panic("APP_ENV not set")
	}

	return appEnv == Production
}
