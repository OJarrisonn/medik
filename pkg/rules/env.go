package rules

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
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

type ValidateEnvOneOfError struct {
	Rule *ValidateEnvOneOfRule
	Value string
}

func (e *ValidateEnvOneOfError) Error() string {
	return fmt.Sprintf("Environment variable %v has value %v which is not one of %v", e.Rule.EnvVar, e.Value, e.Rule.Options)
}

// A rule to check if an environment variable is set and matches one of a list of possible values
type ValidateEnvOneOfRule struct {
	EnvVar string
	Options []string
}

// ValidateEnvOneOf checks if an environment variable is set and matches one of a list of possible values
// TODO: Make faster lookups
func (r *ValidateEnvOneOfRule) Validate() (bool, error) {
	val, ok := os.LookupEnv(r.EnvVar)

	if !ok {
		return false, nil
	}

	for _, opt := range r.Options {
		if val == opt {
			return true, nil
		}
	}

	return false, &ValidateEnvOneOfError{r, val}
}

type ValidateEnvIntegerError struct {
	EnvVar string
	Value string
}

func (e *ValidateEnvIntegerError) Error() string {
	return fmt.Sprintf("Environment variable %v has value %v which is not a number", e.EnvVar, e.Value)
}

// A rule to check if an environment variable is a number
type ValidateEnvIntegerRule struct {
	EnvVar string
}

// ValidateEnvNumber checks if an environment variable is a number
func (r *ValidateEnvIntegerRule) Validate() (bool, error) {
	val, ok := os.LookupEnv(r.EnvVar)

	if !ok {
		return false, nil
	}

	_, err := strconv.Atoi(val)

	if err != nil {
		return false, &ValidateEnvIntegerError{r.EnvVar, val}
	}

	return true, nil
}