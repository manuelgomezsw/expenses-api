package customdate

import (
	"strings"
	"time"
)

const (
	StandardDateTimeFormat = "2006-01-02T15:04:05Z07:00"
)

func SetToNoon(dateStr string) string {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return ""
	}

	t, err := time.Parse(StandardDateTimeFormat, dateStr)
	if err != nil {
		return dateStr
	}

	// Crea la fecha con hora 12:00 local
	noonLocal := time.Date(
		t.Year(), t.Month(), t.Day(),
		12, 0, 0, 0,
		time.Local,
	)

	return noonLocal.Format(StandardDateTimeFormat)
}

func ParseAndFormatDateMySql(dateStr string) (string, error) {
	if dateStr == "" {
		return "", nil
	}

	// Intentar parsear en formato RFC3339 (por ejemplo, "2025-01-24T05:00:00.000Z")
	parsedTime, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return "", err
	}

	// Convertir al layout para MySQL "YYYY-MM-DD HH:MM:SS"
	return parsedTime.Format("2006-01-02 15:04:05"), nil
}
