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
	ErrorLong          = errors.New("Too long")
	ErrorShort         = errors.New("Too short")
	ErrorInvalidChar   = errors.New("Invalid characters")
	ErrorInvalidFormat = errors.New("Invalid format")
	ErrorInvalidLang   = errors.New("Invalid language")
)

func ValidateEmail(in string) error {
	if len(in) > maxEmailLen {
		return ErrorLong
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9.]+@[a-zA-Z0-9.]+\.[a-z]+$`)
	if !re.Match([]byte(in)) {
		return ErrorInvalidFormat
	}
	return nil
}

func ValidatePassword(in string) error {
	switch {
	case len(in) < minPasswordLen:
		return ErrorShort
	case len(in) > maxPasswordLen:
		return ErrorLong
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !re.Match([]byte(in)) {
		return ErrorInvalidChar
	}
	return nil
}

func ValidateLogin(in string) error {
	if len(in) > maxLoginLen {
		return ErrorLong
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9_.]+$`)
	if !re.Match([]byte(in)) {
		return ErrorInvalidChar
	}

	return nil
}

func ValidateLang(in string) error {
	if in != "ru" && in != "en" {
		return ErrorInvalidLang
	}

	return nil
}
