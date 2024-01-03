package api

import (
	"github.com/budiharyonoo/simple_bank/utils"
	"github.com/go-playground/validator/v10"
)

// validCurrency is custom validator for validate supported currency
// list of supported currency available on utils/currency.go
var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return utils.IsSupportedCurrency(currency)
	}

	return false
}
