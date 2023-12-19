package dateutil

import "time"

func GetStartOfDay(t time.Time, location *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, location)
}

func GetEndOfDay(t time.Time, location *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, location)
}

func IsDateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()

	return y1 == y2 && m1 == m2 && d1 == d2
}
