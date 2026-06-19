package httpapi

import (
	"regexp"
	"strconv"
)

var emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func validateStudentID(raw string) (int64, bool) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 || id > 1_000_000_000 {
		return 0, false
	}
	return id, true
}

func validateEmail(email string) bool {
	if len(email) == 0 || len(email) > 254 {
		return false
	}
	return emailPattern.MatchString(email)
}
