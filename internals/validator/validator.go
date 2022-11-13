// Filename : internal/validator/validator.go

package validator

import (
	"regexp"
)

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	PhoneRX = regexp.MustCompile(`^\+?\(?[0-9]{3}\)?\s?-\s?[0-9]{3}\s?-\s?[0-9]{4}$`)
)

// create a type that wraps the validation errors map

type Validator struct {
	Errors map[string]string
}

// Create a new instance of Validator
func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// Valid checks the Errors map for entries
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// In() checks if elements can be found in a provided list of elements
func In(elements string, list ...string) bool {
	for i := range elements {

		if elements == list[i] {
			return true
		}

	}
	return false
}

// Matches() returns true if the provided string matches a specific regexp pattern
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// AddError() adds an error entry to the Error map
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check() preforms the validation checks and calls the AddError method in turn if there is an error
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// Unique() checks that there are no repeating values in the slice
func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)
	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(uniqueValues) == len(values)
}

func ValidateForum(v *validator.Validator, forum *Forum) {
	// Use the Check() method to execute our validation checks
	v.Check(forum.Name != "", "name", "must be provided")
	v.Check(len(forum.Name) <= 200, "name", "must not be more than 200 bytes long")

	v.Check(forum.Message != "", "message", "must be provided")
	v.Check(len(forum.Message) <= 2000, "message", "must not be more than 2000 bytes long")
}
