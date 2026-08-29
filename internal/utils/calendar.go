package utils

// IsLeapYear determines if a given year is a leap year.
func IsLeapYear(year int) bool {
	if year%400 == 0 {
		return true
	}
	return year%4 == 0 && year%100 != 0
}
