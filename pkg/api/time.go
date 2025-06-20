package api

import "time"

func ConvertTimeToStr(v time.Time) string {
	return v.Format(time.RFC3339)
}
