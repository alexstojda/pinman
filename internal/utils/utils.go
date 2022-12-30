package utils

import "time"

func PtrString(str string) *string {
	return &str
}

func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}
