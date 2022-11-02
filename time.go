package goutils

import "time"

var VnLocation *time.Location

// Load default timezone
func LoadLocation() {
	location := Env("TIMEZONE", "Asia/Ho_Chi_Minh")

	loc, err := time.LoadLocation(location)
	if err != nil {
		Panic(err)
	}
	VnLocation = loc
}

// Now() in UTC+7
func Now() time.Time {
	return time.Now().In(VnLocation)
}

// Today() in UTC+7
func Today() time.Time {
	now := Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, VnLocation)
}

// Yesterday() in UTC+7
func Yesterday() time.Time {
	return Today().Add(-24 * time.Hour)
}

// The first date of year in UTC+7
func FirstDateOfYear() time.Time {
	return time.Date(Now().Year(), 1, 1, 0, 0, 0, 0, VnLocation)
}

// Format time to string with RFC3339
func TimeStr(t time.Time) string {
	return t.Format(time.RFC3339)
}

// Parse string to time with RFC3339
func ParseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		Error(err)
		return time.Time{}
	}
	return t
}

// Count days between 2 dates
func CountDays(from time.Time, to time.Time) int {
	return int(to.Sub(from).Hours() / 24)
}

// Count days in a year
func CountDaysInYear(year int) int {
	beginOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, VnLocation)
	endOfYear := time.Date(year, time.December, 31, 0, 0, 0, 0, VnLocation)
	diff := endOfYear.Sub(beginOfYear)
	return int(diff.Hours() / 24)
}
