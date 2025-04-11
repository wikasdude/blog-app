package utils

import (
	"net/http"
	"regexp"
	"strings"
)

func IsValidEmail(email string) bool {

	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func findUSerIDAndRole(w http.ResponseWriter, r *http.Request) (int, string) {
	tokenString := r.Header.Get("Authorization")
	claims, err := ValidateJWT(strings.TrimPrefix(tokenString, "Bearer "))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return 0, ""
	}
	userID := claims.UserID
	role := claims.Role
	return userID, role
}
