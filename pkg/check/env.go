package check

import (
	"fmt"
	"os"
	"regexp"
)

type EnvError struct {
	EnvVar string
	Status string
}

func (e *EnvError) Error() string {
	return e.EnvVar + " is " + e.Status
}

// CheckEnv checks if an environment variable is set
func CheckEnv(name string) bool {
	_, ok := os.LookupEnv(name)

	return ok
}

// ValidateEnvRegex checks if an environment variable is set and matches a regex
func ValidateEnvRegex(name string, regex string) (bool, error) {
	val, ok := os.LookupEnv(name)

	if !ok {
		return false, nil
	}

	regexp, err := regexp.Compile(regex)

	if err != nil {
		return false, err
	}

	if !regexp.MatchString(val) {
		return false, &EnvError{EnvVar: name, Status: fmt.Sprintf("%v which doesn't match %v", val, regex)}
	}

	return true, nil
}