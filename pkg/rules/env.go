package rules

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type EnvIsSetError struct {
	EnvVar string
}

func (e *EnvIsSetError) Error() string {
	return fmt.Sprintf("Environment variable %v is not set", e.EnvVar)
}

// A rule to simply check if a given environment variable is set
// This rule cannot heal the system
type EnvIsSet struct {
	EnvVar string
}

func (r *EnvIsSet) Validate() (bool, error) {
	_, ok := os.LookupEnv(r.EnvVar)

	if !ok {
		return false, &EnvIsSetError{r.EnvVar}
	}

	return true, nil
}

type EnvRegexError struct {
	EnvVar,
	Value,
	Regex string
}

func (e *EnvRegexError) Error() string {
	return fmt.Sprintf("Environment variable %v has value %v which doesn't match %v", e.EnvVar, e.Value, e.Regex)
}

// A rule to check if an environment variable is set and matches a regex
type EnvRegex struct {
	EnvVar string
	Regex  string
}

// ValidateEnvRegex checks if an environment variable is set and matches a regex
func (r *EnvRegex) Validate() (bool, error) {
	val, ok := os.LookupEnv(r.EnvVar)

	if !ok {
		return false, nil
	}

	regexp, rerr := regexp.Compile(r.Regex)

	if rerr != nil {
		return false, rerr
	}

	if !regexp.MatchString(val) {
		return false, &EnvRegexError{r.EnvVar, val, r.Regex}
	}

	return true, nil
}

type EnvOptionError struct {
	Rule  *EnvOption
	Value string
}

func (e *EnvOptionError) Error() string {
	return fmt.Sprintf("Environment variable %v has value %v which is not one of %v", e.Rule.EnvVar, e.Value, e.Rule.Options)
}

// A rule to check if an environment variable is set and matches one of a list of possible values
type EnvOption struct {
	EnvVar  string
	Options []string
}

// ValidateEnvOneOf checks if an environment variable is set and matches one of a list of possible values
// TODO: Make faster lookups
func (r *EnvOption) Validate() (bool, error) {
	val, ok := os.LookupEnv(r.EnvVar)

	if !ok {
		return false, nil
	}

	for _, opt := range r.Options {
		if val == opt {
			return true, nil
		}
	}

	return false, &EnvOptionError{r, val}
}

type EnvIntegerError struct {
	EnvVar string
	Value  string
}

func (e *EnvIntegerError) Error() string {
	return fmt.Sprintf("Environment variable %v has value %v which is not a number", e.EnvVar, e.Value)
}

// A rule to check if an environment variable is a number
type EnvInteger struct {
	EnvVar string
}

// ValidateEnvNumber checks if an environment variable is a number
func (r *EnvInteger) Validate() (bool, error) {
	val, ok := os.LookupEnv(r.EnvVar)

	if !ok {
		return false, nil
	}

	_, err := strconv.Atoi(val)

	if err != nil {
		return false, &EnvIntegerError{r.EnvVar, val}
	}

	return true, nil
}


type EnvFloatError struct {
	EnvVar string
	Value  string
}

func (e *EnvFloatError) Error() string {
	return fmt.Sprintf("Environment variable %v has value %v which is not a float", e.EnvVar, e.Value)
}

// A rule to check if an environment variable is a number
type EnvFloat struct {
	EnvVar string
}

// ValidateEnvNumber checks if an environment variable is a number
func (r *EnvFloat) Validate() (bool, error) {
	val, ok := os.LookupEnv(r.EnvVar)

	if !ok {
		return false, nil
	}

	_, err := strconv.ParseFloat(val, 64)

	if err != nil {
		return false, &EnvFloatError{r.EnvVar, val}
	}

	return true, nil
}