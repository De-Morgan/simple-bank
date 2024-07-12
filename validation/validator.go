package validation

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func validateStringLen(value string, min, max int) error {
	if n := len(value); n < min || n > max {
		return fmt.Errorf("must be between %d-%d characters", min, max)
	}
	return nil
}

func ValidateUsername(username string) (err error) {
	err = validateStringLen(username, 3, 100)
	if err != nil {
		return
	}
	if !isValidUsername(username) {
		return fmt.Errorf("it must contain only lowercase letter digit or underscore")
	}
	return
}

func ValidatePassword(password string) (err error) {
	err = validateStringLen(password, 6, 100)
	return
}

func ValidEmailAddress(email string) (err error) {
	if _, err = mail.ParseAddress(email); err != nil {
		err = fmt.Errorf("is not a valid email address")
	}
	return
}

func ValidateFullName(name string) (err error) {
	err = validateStringLen(name, 3, 100)
	if err != nil {
		return
	}
	if !isValidFullName(name) {
		return fmt.Errorf("it must contain only letter and spaces")
	}
	return
}
