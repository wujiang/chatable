package forms

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	passwordMinLen = 10

	emailRE = regexp.MustCompile("^[\\w!#$%&'*+\\-/=?^`{|}~.]+@(?:[a-z0-9][a-z0-9\\-]{0,62}\\.)+(?:[a-z]{2,63}|xn--[a-z0-9\\-]{2,59})$")
	phoneRE = regexp.MustCompile("^[0-9]+$")

	ErrEmailFormat   = errors.New("please enter a valid email address")
	ErrPasswordLen   = errors.New(fmt.Sprintf("please enter a password with at least %d characters", passwordMinLen))
	ErrPasswordSpace = errors.New("please enter a password without any space")
	ErrPhoneFormat   = errors.New("please enter a valid phone number")
	ErrNameFormat    = errors.New("please enter a valid name")
	ErrStringValue   = errors.New("invalid string value")
	ErrIntValue      = errors.New("invalid int value")
)

type ValidatorFunc func(interface{}) error

func ReValidator(pattern string, errorMsg string) ValidatorFunc {
	validRe := regexp.MustCompile(pattern)
	return ValidatorFunc(func(value interface{}) error {
		strValue, ok := value.(string)
		// check type assertion
		if !ok {
			return ErrStringValue
		}
		// check regular expression
		if !validRe.MatchString(strValue) {
			return errors.New(errorMsg)
		}
		return nil
	})
}

func RangeValidator(min, max int, errorMsg string) ValidatorFunc {
	return ValidatorFunc(func(value interface{}) error {
		intValue, ok := value.(int)
		// check type assertion
		if !ok {
			return ErrIntValue
		}
		// check range
		if intValue < min || intValue > max {
			return errors.New(errorMsg)
		}
		return nil
	})
}

func EmailValidator() ValidatorFunc {
	return ValidatorFunc(func(value interface{}) error {
		strValue, ok := value.(string)
		if !ok {
			return ErrEmailFormat
		}
		if !emailRE.MatchString(strValue) {
			return ErrEmailFormat
		}
		return nil
	})
}

func PasswordValidator() ValidatorFunc {
	return ValidatorFunc(func(value interface{}) error {
		strValue, ok := value.(string)
		if !ok {
			return ErrPasswordLen
		}
		strValue = strings.TrimSpace(strValue)
		if strings.ContainsRune(strValue, ' ') {
			return ErrPasswordSpace
		}
		if len(strValue) < passwordMinLen {
			return ErrPasswordLen
		}
		return nil
	})
}

func PhoneNumberValidator() ValidatorFunc {
	return ValidatorFunc(func(value interface{}) error {
		strValue, ok := value.(string)
		if !ok {
			return ErrPhoneFormat
		}
		strValue = strings.TrimSpace(strValue)
		if !phoneRE.MatchString(strValue) {
			return ErrPhoneFormat
		}
		return nil
	})
}

func NameValidator() ValidatorFunc {
	return ValidatorFunc(func(value interface{}) error {
		strValue, ok := value.(string)
		if !ok {
			return ErrPhoneFormat
		}
		strValue = strings.TrimSpace(strValue)
		// if !nameRE.MatchString(strValue) {
		// 	return ErrNameFormat
		// }
		return nil
	})
}
