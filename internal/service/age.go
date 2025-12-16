package service

import "time"

// CalculateAge returns full years between dob and now.
func CalculateAge(dob time.Time, now time.Time) int {
	age := now.Year() - dob.Year()
	// if birthday hasn't happened yet this year, subtract 1
	if now.Month() < dob.Month() || (now.Month() == dob.Month() && now.Day() < dob.Day()) {
		age--
	}
	if age < 0 { // guard for future dates
		return 0
	}
	return age
}
