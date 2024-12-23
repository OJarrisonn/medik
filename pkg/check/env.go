package check

import (
	"fmt"
	"os"
	"regexp"
)

// CheckEnv checks if an environment variable is set
func CheckEnv(name string) bool {
	_, ok := os.LookupEnv(name)

	return ok
}

type ValidateEnvRegexError struct {
	EnvVar,
	Value,
	Regex string
}

func (e ValidateEnvRegexError) Error() string {
	return fmt.Sprintf("Environment variable %v has value %v which doesn't match %v", e.EnvVar, e.Value, e.Regex)
}

// ValidateEnvRegex checks if an environment variable is set and matches a regex
func ValidateEnvRegex(name string, regex string) (bool, error) {
	val, ok := os.LookupEnv(name)

	if !ok {
		return false, nil
	}

	regexp, rerr := regexp.Compile(regex)

	if rerr != nil {
		return false, rerr
	}

	if !regexp.MatchString(val) {
		return false, ValidateEnvRegexError{name, val, regex}
	}

	return true, nil
}
