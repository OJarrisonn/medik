package exams

import (
	"fmt"
	neturl "net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// A rule to simply check if a given environment variable is set
// This rule cannot heal the system
type EnvIsSet struct {
	EnvVar string
}

func (r *EnvIsSet) Examinate() (bool, error) {
	_, ok := os.LookupEnv(r.EnvVar)

	if !ok {
		return false, fmt.Errorf("environment variable %v is not set", r.EnvVar)
	}

	return true, nil
}

type EnvIsSetNotEmpty struct {
	EnvVar string
}

func (r *EnvIsSetNotEmpty) Examinate() (bool, error) {
	val, ok := os.LookupEnv(r.EnvVar)

	if !ok {
		return false, fmt.Errorf("environment variable %v is not set", r.EnvVar)
	}

	if strings.TrimSpace(val) == "" {
		return false, fmt.Errorf("environment variable %v is set to an empty string \"%v\"", r.EnvVar, val)
	}

	return true, nil
}

// A rule to check if an environment variable is set and matches a regex
type EnvRegex struct {
	EnvVar string
	Regex  string
}

// ValidateEnvRegex checks if an environment variable is set and matches a regex
func (r *EnvRegex) Examinate() (bool, error) {
	val, ok := os.LookupEnv(r.EnvVar)

	if !ok {
		return false, nil
	}

	regexp, rerr := regexp.Compile(r.Regex)

	if rerr != nil {
		return false, rerr
	}

	if !regexp.MatchString(val) {
		return false, fmt.Errorf("environment variable %v has value %v which doesn't match %v", r.EnvVar, val, r.Regex)
	}

	return true, nil
}

// A rule to check if an environment variable is set and matches one of a list of possible values
type EnvOption struct {
	EnvVar  string
	Options []string
}

// ValidateEnvOneOf checks if an environment variable is set and matches one of a list of possible values
// TODO: Make faster lookups
func (r *EnvOption) Examinate() (bool, error) {
	val, ok := os.LookupEnv(r.EnvVar)

	if !ok {
		return false, nil
	}

	for _, opt := range r.Options {
		if val == opt {
			return true, nil
		}
	}

	return false, fmt.Errorf("environment variable %v has value %v which is not one of %v", r.EnvVar, val, r.Options)
}

// A rule to check if an environment variable is a number
type EnvInteger struct {
	EnvVar string
}

// ValidateEnvNumber checks if an environment variable is a number
func (r *EnvInteger) Examinate() (bool, error) {
	val, ok := os.LookupEnv(r.EnvVar)

	if !ok {
		return false, nil
	}

	_, err := strconv.Atoi(val)

	if err != nil {
		return false, fmt.Errorf("environment variable %v has value %v which is not a number", r.EnvVar, val)
	}

	return true, nil
}

// A rule to check if an environment variable is a number within a range
// Min and Max are inclusive
type EnvIntegerRange struct {
	EnvVar string
	Min    int
	Max    int
}

// ValidateEnvNumber checks if an environment variable is a number
func (r *EnvIntegerRange) Examinate() (bool, error) {
	val, ok := os.LookupEnv(r.EnvVar)

	if !ok {
		return false, nil
	}

	num, err := strconv.Atoi(val)

	if err != nil {
		return false, fmt.Errorf("environment variable %v has value %v which is not a number", r.EnvVar, val)
	}

	if num < r.Min || num > r.Max {
		return false, fmt.Errorf("environment variable %v has value %v which is not in the int range [%v,%v]", r.EnvVar, val, r.Min, r.Max)
	}

	return true, nil
}

// A rule to check if an environment variable is a number
type EnvFloat struct {
	EnvVar string
}

// ValidateEnvNumber checks if an environment variable is a number
func (r *EnvFloat) Examinate() (bool, error) {
	val, ok := os.LookupEnv(r.EnvVar)

	if !ok {
		return false, nil
	}

	_, err := strconv.ParseFloat(val, 64)

	if err != nil {
		return false, fmt.Errorf("environment variable %v has value %v which is not a float", r.EnvVar, val)
	}

	return true, nil
}

// A rule to check if an environment variable is a number within a range
// Min and Max are inclusive
type EnvFloatRange struct {
	EnvVar string
	Min    float64
	Max    float64
}

// ValidateEnvNumber checks if an environment variable is a number
func (r *EnvFloatRange) Examinate() (bool, error) {
	val, ok := os.LookupEnv(r.EnvVar)

	if !ok {
		return false, nil
	}

	num, err := strconv.ParseFloat(val, 64)

	if err != nil {
		return false, fmt.Errorf("environment variable %v has value %v which is not a float", r.EnvVar, val)
	}

	if num < r.Min || num > r.Max {
		return false, fmt.Errorf("environment variable %v has value %v which is not in the float range [%v,%v]", r.EnvVar, val, r.Min, r.Max)
	}

	return true, nil
}

// A rule to check if an environment variable is set to a file that exists
type EnvFile struct {
	EnvVar string
}

// ValidateEnvFile checks if an environment variable is set to a file that exists
func (r *EnvFile) Examinate() (bool, error) {
	val, ok := os.LookupEnv(r.EnvVar)

	if !ok {
		return false, nil
	}

	_, err := os.Stat(val)

	if err != nil {
		return false, fmt.Errorf("environment variable %v has value %v which is not a file, %v", r.EnvVar, val, err)
	}

	return true, nil
}

// A rule to check if an environment variable is set to a directory that exists
type EnvDir struct {
	EnvVar string
}

// ValidateEnvDir checks if an environment variable is set to a directory that exists
func (r *EnvDir) Examinate() (bool, error) {
	val, ok := os.LookupEnv(r.EnvVar)

	if !ok {
		return false, nil
	}

	stat, err := os.Stat(val)

	if err != nil {
		return false, fmt.Errorf("environment variable %v has value %v which is not a directory, %v", r.EnvVar, val, err)
	}

	if !stat.IsDir() {
		return false, fmt.Errorf("environment variable %v has value %v which is not a directory", r.EnvVar, val)
	}

	return true, nil
}

// A rule to check if an environment variable is set to an IP address
type EnvIpv4Addr struct {
	EnvVar string
}

// ValidateEnvIpAddr checks if an environment variable is set to an IP address
func (r *EnvIpv4Addr) Examinate() (bool, error) {
	val, ok := os.LookupEnv(r.EnvVar)

	if !ok {
		return false, nil
	}

	regexp := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)

	if !regexp.MatchString(val) {
		return false, fmt.Errorf("environment variable %v has value %v which is not an IPv4 address", r.EnvVar, val)
	}

	return true, nil
}

// A rule to check if an environment variable is set to an IP address
type EnvIpv6Addr struct {
	EnvVar string
}

// ValidateEnvIpAddr checks if an environment variable is set to an IP address
func (r *EnvIpv6Addr) Examinate() (bool, error) {
	val, ok := os.LookupEnv(r.EnvVar)

	if !ok {
		return false, nil
	}

	regexp := regexp.MustCompile(`^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`)

	if !regexp.MatchString(val) {
		return false, fmt.Errorf("environment variable %v has value %v which is not an IPv6 address", r.EnvVar, val)
	}

	return true, nil
}

// A rule to check if an environment variable is set to an IP address
type EnvIpAddr struct {
	EnvVar string
}

// ValidateEnvIpAddr checks if an environment variable is set to an IP address (v4 or v6)
func (r *EnvIpAddr) Examinate() (bool, error) {
	ipv4 := &EnvIpv4Addr{r.EnvVar}

	if ok, _ := ipv4.Examinate(); ok {
		return true, nil
	}

	ipv6 := &EnvIpv6Addr{r.EnvVar}

	if ok, _ := ipv6.Examinate(); ok {
		return true, nil
	}

	return false, fmt.Errorf("environment variable %v has value %v which is not an IP address", r.EnvVar, os.Getenv(r.EnvVar))
}

// A rule to check if an environment variable is set to a hostname
type EnvHostname struct {
	EnvVar           string
	Protocol         string
}

// ValidateEnvHostname checks if an environment variable is set to a hostname
func (r *EnvHostname) Examinate() (bool, error) {
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

	if r.Protocol == "" {
		return true, nil
	}

	if url.Scheme != r.Protocol {
		return false, fmt.Errorf("environment variable %v has value %v which is not a valid URL with protocol %v", r.EnvVar, rawUrl, r.Protocol)
	}

	return true, nil
}
