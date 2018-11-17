package validator

import (
	"errors"
	"regexp"
)

const (
	maxEmailLen    = 50
	maxLoginLen    = 30
	maxPasswordLen = 50
	minPasswordLen = 8
)

var (
	errorLong          = errors.New("Too long")
	errorShort         = errors.New("Too short")
	errorInvalidChar   = errors.New("Invalid characters")
	errorInvalidFormat = errors.New("Invalid format")
	errorInvalidLang   = errors.New("Invalid language")
)

func ValidateEmail(in string) error {
	if len(in) > maxEmailLen {
		return errorLong
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9.]+@[a-zA-Z0-9.]+\.[a-z]+$`)
	if !re.Match([]byte(in)) {
		return errorInvalidFormat
	}
	return nil
}

func ValidatePassword(in string) error {
	switch {
	case len(in) < minPasswordLen:
		return errorShort
	case len(in) > maxPasswordLen:
		return errorLong
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !re.Match([]byte(in)) {
		return errorInvalidChar
	}
	return nil
}

func ValidateLogin(in string) error {
	if len(in) > maxLoginLen {
		return errorLong
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9_.]+$`)
	if !re.Match([]byte(in)) {
		return errorInvalidChar
	}

	return nil
}

func ValidateLang(in string) error {
	if in != "ru" && in != "en" {
		return errorInvalidLang
	}

	return nil
}
