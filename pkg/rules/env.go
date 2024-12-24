package rules

import (
	"fmt"
	neturl "net/url"
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


type EnvIntegerRangeError struct {
	EnvVar string
	Value  string
	Min    int
	Max    int
}

func (e *EnvIntegerRangeError) Error() string {
	return fmt.Sprintf("Environment variable %v has value %v which is not in the int range [%v,%v)", e.EnvVar, e.Value, e.Min, e.Max)
}

// A rule to check if an environment variable is a number within a range
// Min is inclusive, Max is exclusive
type EnvIntegerRange struct {
	EnvVar string
	Min    int
	Max    int
}

// ValidateEnvNumber checks if an environment variable is a number
func (r *EnvIntegerRange) Validate() (bool, error) {
	val, ok := os.LookupEnv(r.EnvVar)

	if !ok {
		return false, nil
	}

	num, err := strconv.Atoi(val)

	if err != nil {
		return false, &EnvIntegerError{r.EnvVar, val}
	}

	if num < r.Min || num >= r.Max {
		return false, &EnvIntegerRangeError{r.EnvVar, val, r.Min, r.Max}
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

type EnvFloatRangeError struct {
	EnvVar string
	Value  string
	Min    float64
	Max    float64
}

func (e *EnvFloatRangeError) Error() string {
	return fmt.Sprintf("Environment variable %v has value %v which is not in the float range [%v,%v]", e.EnvVar, e.Value, e.Min, e.Max)
}

// A rule to check if an environment variable is a number within a range
// Min and Max are inclusive
type EnvFloatRange struct {
	EnvVar string
	Min    float64
	Max    float64
}

// ValidateEnvNumber checks if an environment variable is a number
func (r *EnvFloatRange) Validate() (bool, error) {
	val, ok := os.LookupEnv(r.EnvVar)

	if !ok {
		return false, nil
	}

	num, err := strconv.ParseFloat(val, 64)

	if err != nil {
		return false, &EnvFloatError{r.EnvVar, val}
	}

	if num < r.Min || num > r.Max {
		return false, &EnvFloatRangeError{r.EnvVar, val, r.Min, r.Max}
	}

	return true, nil
}

type EnvFileError struct {
	EnvVar string
	Value  string
}

func (e *EnvFileError) Error() string {
	return fmt.Sprintf("Environment variable %v has value %v which is not a file", e.EnvVar, e.Value)
}

// A rule to check if an environment variable is set to a file that exists
type EnvFile struct {
	EnvVar string
}

// ValidateEnvFile checks if an environment variable is set to a file that exists
func (r *EnvFile) Validate() (bool, error) {
	val, ok := os.LookupEnv(r.EnvVar)

	if !ok {
		return false, nil
	}

	_, err := os.Stat(val)

	if err != nil {
		return false, &EnvFileError{r.EnvVar, val}
	}

	return true, nil
}

type EnvDirError struct {
	EnvVar string
	Value  string
}

func (e *EnvDirError) Error() string {
	return fmt.Sprintf("Environment variable %v has value %v which is not a directory", e.EnvVar, e.Value)
}

// A rule to check if an environment variable is set to a directory that exists
type EnvDir struct {
	EnvVar string
}

// ValidateEnvDir checks if an environment variable is set to a directory that exists
func (r *EnvDir) Validate() (bool, error) {
	val, ok := os.LookupEnv(r.EnvVar)

	if !ok {
		return false, nil
	}

	stat, err := os.Stat(val)

	if err != nil {
		return false, &EnvDirError{r.EnvVar, val}
	}

	if !stat.IsDir() {
		return false, &EnvDirError{r.EnvVar, val}
	}

	return true, nil
}

type EnvIpv4AddrError struct {
	EnvVar string
	Value  string
}

func (e *EnvIpv4AddrError) Error() string {
	return fmt.Sprintf("Environment variable %v has value %v which is not an IP address", e.EnvVar, e.Value)
}

// A rule to check if an environment variable is set to an IP address
type EnvIpv4Addr struct {
	EnvVar string
}

// ValidateEnvIpAddr checks if an environment variable is set to an IP address
func (r *EnvIpv4Addr) Validate() (bool, error) {
	val, ok := os.LookupEnv(r.EnvVar)

	if !ok {
		return false, nil
	}

	regexp := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	
	if !regexp.MatchString(val) {
		return false, &EnvIpv4AddrError{r.EnvVar, val}
	}

	return true, nil
}

type EnvIpv6AddrError struct {
	EnvVar string
	Value  string
}

func (e *EnvIpv6AddrError) Error() string {
	return fmt.Sprintf("Environment variable %v has value %v which is not an IP address", e.EnvVar, e.Value)
}

// A rule to check if an environment variable is set to an IP address
type EnvIpv6Addr struct {
	EnvVar string
}

// ValidateEnvIpAddr checks if an environment variable is set to an IP address
func (r *EnvIpv6Addr) Validate() (bool, error) {
	val, ok := os.LookupEnv(r.EnvVar)

	if !ok {
		return false, nil
	}

	regexp := regexp.MustCompile(`^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`)

	if !regexp.MatchString(val) {
		return false, &EnvIpv6AddrError{r.EnvVar, val}
	}

	return true, nil
}

type EnvIpAddrError struct {
	EnvVar string
	Value  string
}

func (e *EnvIpAddrError) Error() string {
	return fmt.Sprintf("Environment variable %v has value %v which is not an IP address", e.EnvVar, e.Value)
}

// A rule to check if an environment variable is set to an IP address
type EnvIpAddr struct {
	EnvVar string
}

// ValidateEnvIpAddr checks if an environment variable is set to an IP address (v4 or v6)
func (r *EnvIpAddr) Validate() (bool, error) {
	ipv4 := &EnvIpv4Addr{r.EnvVar}
	
	if ok, _ := ipv4.Validate(); ok {
		return true, nil
	}
	
	ipv6 := &EnvIpv6Addr{r.EnvVar}

	if ok, _ := ipv6.Validate(); ok {
		return true, nil
	}

	return false, &EnvIpAddrError{r.EnvVar, os.Getenv(r.EnvVar)}
}

type EnvHostnameInvalidProtoError struct {
	Rule *EnvHostname
	Value  string
}

func (e *EnvHostnameInvalidProtoError) Error() string {
	return fmt.Sprintf("Environment variable %v has value %v which is not a valid URL with protocol %v", e.Rule.EnvVar, e.Value, e.Rule.Protocol)
}

// A rule to check if an environment variable is set to a hostname
type EnvHostname struct {
	EnvVar string
	SpecificProtocol bool
	Protocol string	
}

// ValidateEnvHostname checks if an environment variable is set to a hostname
func (r *EnvHostname) Validate() (bool, error) {
	val, ok := os.LookupEnv(r.EnvVar)

	if !ok {
		return false, nil
	}

	// Check if the hostname is a valid URL
	return r.validateUrl(val)
}

func (r *EnvHostname) validateUrl(rawUrl string) (bool, error) {
	url, err := neturl.Parse(rawUrl)

	if err != nil {
		return false, err
	}

	if !r.SpecificProtocol {
		return true, nil
	}

	if url.Scheme != r.Protocol {
		return false, &EnvHostnameInvalidProtoError{r, rawUrl}
	}

	return true, nil
}
