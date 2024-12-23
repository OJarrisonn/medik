package rules

import (
	"fmt"
	"os"
	"regexp"
)

type CheckEnvError struct {
	EnvVar string
}

func (e *CheckEnvError) Error() string {
	return fmt.Sprintf("Environment variable %v is not set", e.EnvVar)
}

// A rule to simply check if a given environment variable is set
// This rule cannot heal the system
type CheckEnvRule struct {
	EnvVar string
}

func (r *CheckEnvRule) Validate() (bool, error) {
	_, ok := os.LookupEnv(r.EnvVar)

	if !ok {
		return false, &CheckEnvError{r.EnvVar}
	}

	return true, nil
}

type ValidateEnvRegexError struct {
	EnvVar,
	Value,
	Regex string
}

func (e *ValidateEnvRegexError) Error() string {
	return fmt.Sprintf("Environment variable %v has value %v which doesn't match %v", e.EnvVar, e.Value, e.Regex)
}

// A rule to check if an environment variable is set and matches a regex
type ValidateEnvRegexRule struct {
	EnvVar string
	Regex  string
}

// ValidateEnvRegex checks if an environment variable is set and matches a regex
func (r *ValidateEnvRegexRule) Validate() (bool, error) {
	val, ok := os.LookupEnv(r.EnvVar)

	if !ok {
		return false, nil
	}

	regexp, rerr := regexp.Compile(r.Regex)

	if rerr != nil {
		return false, rerr
	}

	if !regexp.MatchString(val) {
		return false, &ValidateEnvRegexError{r.EnvVar, val, r.Regex}
	}

	return true, nil
}
