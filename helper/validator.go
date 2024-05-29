package helper

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
	"strings"
	"unicode"
)

func IsAlphaNumeric(s string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(s)
}

func IsASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func IsContainWhiteSpace(s string) bool {
	return regexp.MustCompile(`\s`).MatchString(s)
}

func ValidatePassword(password string) (bool, string) {

	var numeric, lowercase, uppercase, specialChar bool

	if len(password) < 8 {
		return false, "password length can't less than 8"
	}

	numeric = regexp.MustCompile(`\d`).MatchString(password)
	lowercase = regexp.MustCompile(`[a-z]`).MatchString(password)
	uppercase = regexp.MustCompile(`[A-Z]`).MatchString(password)
	//!@#$%^&*
	specialChar = strings.ContainsAny(password, "!@#$%^&*")

	if !numeric {
		return false, "password must contain at least one number"
	}

	if !lowercase {
		return false, "password must contain at least one lowercase letter"
	}

	if !uppercase {
		return false, "password must contain at least one uppercase letter"
	}

	if !specialChar {
		return false, "password must contain at least one special char '!@#$%^&*'"
	}

	return true, ""
}

func ValidateStruct(object interface{}) *string {
	validate := validator.New()
	err := validate.Struct(object)
	if err == nil {
		return nil
	}

	var validationMsg *string

	if errv, ok := err.(validator.ValidationErrors); ok {
		msg := fmt.Sprintf("%s - %s", errv[0].Field(), msgForTag(errv[0].Tag()))
		validationMsg = &msg
	}

	return validationMsg
}

func msgForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	}
	return ""
}
