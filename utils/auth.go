package utils

import (
	"regexp"
	"strings"
)

func IsValidEmail(email string) bool {
	pattern := `^[a-z]+(?:\.[a-z]+)?@remotestate\.com$`

	re := regexp.MustCompile(pattern)
	return re.MatchString(email)
}

func ExtractNameFromEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 0 {
		return parts[0] // everything before @
	}
	return ""
}

func IsValidRole(role string) bool {
	roles := []string{"admin", "employee_manager", "asset_manager", "employee"}

	for _, r := range roles {
		if role == r {
			return true
		}
	}
	return false
}

func IsValidUserType(userType string) bool {
	userTypes := []string{"full_time", "intern", "freelancer"}
	for _, ut := range userTypes {
		if userType == ut {
			return true
		}
	}
	return false
}

func IsValidPhone(phone string) bool {
	if len(phone) == 10 {
		return true
	}
	return false
}
