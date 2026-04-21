package validation

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func isUppercase(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return value == value[:]
}

func isPureString(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	re := regexp.MustCompile(CheckOnlyAlphabet)
	return re.MatchString(value)
}

func isDigit(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	re := regexp.MustCompile(Digits)
	return re.MatchString(value)
}

func isNPWP(fl validator.FieldLevel) bool {
	npwp := fl.Field().String()
	re := regexp.MustCompile(NPWP)
	if !re.MatchString(npwp) {
		re := regexp.MustCompile(KTP)
		return re.MatchString(npwp)
	}
	return true
}

func validateBankName(fl validator.FieldLevel) bool {
	bankName := fl.Field().String()

	// Check if the bank name is in the list of valid bank names
	for _, validBank := range BankChoices {
		if bankName == validBank {
			return true
		}
	}
	return false
}

func validateEnum(fl validator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(Enum); ok {
		return value.IsValidEnum()
	}
	return false
}
