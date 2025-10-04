package handler

import (
	"errors"
	"time"
)

// getPreviousMonth calcula el mes anterior en formato YYYY-MM
// Maneja correctamente el cambio de año (ej: 2024-01 → 2023-12)
func getPreviousMonth(month string) (string, error) {
	// Validar y parsear el mes
	date, err := time.Parse("2006-01", month)
	if err != nil {
		return "", errors.New("invalid month format, must be YYYY-MM")
	}

	// Restar un mes
	previousMonth := date.AddDate(0, -1, 0)

	// Retornar en formato YYYY-MM
	return previousMonth.Format("2006-01"), nil
}
