package utils

import (
	"errors"
	"regexp"
)

// ValidateLoginPassword validate login & password
func ValidateLoginPassword(login string, password string) error {
	if login == "" {
		return errors.New("Login empty")
	}
	if password == "" {
		return errors.New("Password empty")
	}
	return nil
}

// IsOnlyDigits only digits in string validater
func IsOnlyDigits(s string) bool {
	result, err := regexp.MatchString("^[0-9]*$", s)
	if err != nil {
		return false
	}
	return result
}
