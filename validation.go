package goodidea

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

var validate *validator.Validate

// Register the custom validator
func init() {
	validate = validator.New()
	validate.RegisterValidation("u20", U20Validator)
}

// U20Validator checks if the string  matches a userID which is an xid string
// with a u prefixed. i.e. starts with 'u' followed by 20 alphanumeric characters (a-v, 0-9)
func U20Validator(fl validator.FieldLevel) bool {
	pattern := "^u[a-v0-9]{20}$"
	re := regexp.MustCompile(pattern)
	value := fl.Field().String()
	return re.MatchString(value)
}
