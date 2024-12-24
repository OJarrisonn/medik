package exams

import (
	"fmt"
	neturl "net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Check if a given environment variable is set
// 
// type: env.is-set,
// vars: []string
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

// Check if an environment variable is set and not empty
// Strings with only whitespace are considered empty
//
// type: env.not-empty
// vars: []string
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

// Check if an environment variable is set and matches a regex
// The regex is a string that will be compiled into a regular expression using Go's regexp package
//
// type: env.regex,
// vars: []string,
// regex: string
type EnvRegex struct {
	EnvVar string
	Regex  string
}

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

// Check if an environment variable is set and matches one of a list of possible values
//
// type: env.options,
// vars: []string,
// options: []string
type EnvOption struct {
	EnvVar  string
	Options []string
}

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

// Check if an environment variable is a number
// This rule will check if the environment variable is a number
// 
// type: env.int,
// vars: []string
type EnvInteger struct {
	EnvVar string
}

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

// Check if an environment variable is a number within a range
// Min and Max are inclusive
//
// type: env.int-range,
// vars: []string,
// min: int,
// max: int
type EnvIntegerRange struct {
	EnvVar string
	Min    int
	Max    int
}

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

// Check if an environment variable is a floating number
//
// type: env.float,
// vars: []string
type EnvFloat struct {
	EnvVar string
}

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

// Check if an environment variable is a number within a range
// Min and Max are inclusive
//
// type: env.float-range,
// vars: []string,
// min: float,
// max: float
type EnvFloatRange struct {
	EnvVar string
	Min    float64
	Max    float64
}

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

// TODO: Add support to point to files that don't exist
// Check if an environment variable is set to a file that exists
// 
// type: env.file,
// vars: []string
type EnvFile struct {
	EnvVar string
}

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

// Check if an environment variable is set to a directory that exists
//
// type: env.dir,
// vars: []string
type EnvDir struct {
	EnvVar string
}

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

// Check if an environment variable is set to an IP address
//
// type: env.ipv4,
// vars: []string
type EnvIpv4Addr struct {
	EnvVar string
}

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

// Check if an environment variable is set to an IP address
//
// type: env.ipv6,
// vars: []string
type EnvIpv6Addr struct {
	EnvVar string
}

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

// Check if an environment variable is set to an IP address
//
// type: env.ip,
// vars: []string
type EnvIpAddr struct {
	EnvVar string
}

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
//
// type: env.hostname,
// vars: []string,
// protocol: string
type EnvHostname struct {
	EnvVar           string
	Protocol         string
}

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
