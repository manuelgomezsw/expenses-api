package customdate

func RemoveTime(date string) string {
	if len(date) >= 10 {
		return date[:10]
	}

	return date
}
