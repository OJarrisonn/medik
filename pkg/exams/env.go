package exams

import (
	"fmt"
	neturl "net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/OJarrisonn/medik/pkg/config"
)

// Function to get a parser for a given type `env.*`
// The type is the part after `env.` in the exam type
// Returns the parser and a boolean indicating if the parser was found
func Parser(ty string) (ExamParser, bool) {
	if parser, ok := parsers[ty]; ok {
		return parser, ok
	}

	return nil, false
}

var parsers = map[string]ExamParser{
	"is-set":      getParser[*EnvIsSet](),
	"not-empty":   getParser[*EnvIsSetNotEmpty](),
	"regex":       getParser[*EnvRegex](),
	"options":     getParser[*EnvOption](),
	"int":         getParser[*EnvInteger](),
	"int-range":   getParser[*EnvIntegerRange](),
	"float":       getParser[*EnvFloat](),
	"float-range": getParser[*EnvFloatRange](),
	"file":        getParser[*EnvFile](),
	"dir":         getParser[*EnvDir](),
	"ipv4":        getParser[*EnvIpv4Addr](),
	"ipv6":        getParser[*EnvIpv6Addr](),
	"ip":          getParser[*EnvIpAddr](),
	"hostname":    getParser[*EnvHostname](),
}

func getParser[E Exam]() ExamParser {
	var e E
	return e.Parse
}

// Check if a given environment variable is set
//
// type: env.is-set,
// vars: []string
type EnvIsSet struct {
	Vars []string
}

func (r *EnvIsSet) Type() string {
	return "env.is-set"
}

func (r *EnvIsSet) Parse(config config.Exam) (Exam, error) {
	if config.Type != r.Type() {
		return nil, fmt.Errorf("invalid type %v for env.is-set", config.Type)
	}

	if len(config.Vars) == 0 {
		return nil, fmt.Errorf("vars is not set for env.is-set")
	}

	return &EnvIsSet{config.Vars}, nil
}

func (r *EnvIsSet) Examinate() (bool, error) {
	unset := []string{}

	for _, v := range r.Vars {
		if _, ok := os.LookupEnv(v); !ok {
			unset = append(unset, v)
		}
	}

	if len(unset) > 0 {
		return false, fmt.Errorf("environment variables not set %v", unset)
	}

	return true, nil
}

// Check if an environment variable is set and not empty
// Strings with only whitespace are considered empty
//
// type: env.not-empty
// vars: []string
type EnvIsSetNotEmpty struct {
	Vars []string
}

func (r *EnvIsSetNotEmpty) Type() string {
	return "env.not-empty"
}

func (r *EnvIsSetNotEmpty) Parse(config config.Exam) (Exam, error) {
	if config.Type != r.Type() {
		return nil, fmt.Errorf("invalid type %v for env.not-empty", config.Type)
	}

	if len(config.Vars) == 0 {
		return nil, fmt.Errorf("vars is not set for env.not-empty")
	}

	return &EnvIsSetNotEmpty{config.Vars}, nil
}

func (r *EnvIsSetNotEmpty) Examinate() (bool, error) {
	unset := []string{}
	invalid := []string{}

	for _, v := range r.Vars {
		if val, ok := os.LookupEnv(v); !ok {
			unset = append(unset, v)
		} else if strings.TrimSpace(val) == "" {
			invalid = append(invalid, v)
		}
	}

	err := ""

	if len(unset) > 0 {
		err += fmt.Sprintf("environment variables not set %v\n", unset)
	}

	if len(invalid) > 0 {
		err += fmt.Sprintf("environment variables set to empty strings %v\n", invalid)
	}

	if err != "" {
		return false, fmt.Errorf("%v", err)
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
	Vars  []string
	Regex *regexp.Regexp
}

func (r *EnvRegex) Type() string {
	return "env.regex"
}

func (r *EnvRegex) Parse(config config.Exam) (Exam, error) {
	if config.Type != r.Type() {
		return nil, fmt.Errorf("invalid type %v for env.regex", config.Type)
	}

	if len(config.Vars) == 0 {
		return nil, fmt.Errorf("vars is not set for env.regex")
	}

	if config.Regex == "" {
		return nil, fmt.Errorf("regex is not set for env.regex")
	}

	regexp, rerr := regexp.Compile(config.Regex)

	if rerr != nil {
		return nil, fmt.Errorf("invalid regex %v for env.regex, %v", config.Regex, rerr)
	}

	return &EnvRegex{config.Vars, regexp}, nil
}

func (r *EnvRegex) Examinate() (bool, error) {
	unset := []string{}
	invalid := []string{}

	for _, v := range r.Vars {
		if val, ok := os.LookupEnv(v); !ok {
			unset = append(unset, v)
		} else if !r.Regex.MatchString(val) {
			invalid = append(invalid, v)
		}
	}

	err := ""

	if len(unset) > 0 {
		err += fmt.Sprintf("environment variables not set %v\n", unset)
	}

	if len(invalid) > 0 {
		err += fmt.Sprintf("environment variables not matching regex %v\n", invalid)
	}

	if err != "" {
		return false, fmt.Errorf("%v", err)
	}

	return true, nil
}

// Check if an environment variable is set and matches one of a list of possible values
//
// type: env.options,
// vars: []string,
// options: []string
type EnvOption struct {
	Vars    []string
	Options map[string]bool
}

func (r *EnvOption) Type() string {
	return "env.options"
}

func (r *EnvOption) Parse(config config.Exam) (Exam, error) {
	if config.Type != r.Type() {
		return nil, fmt.Errorf("invalid type %v for env.options", config.Type)
	}

	if len(config.Vars) == 0 {
		return nil, fmt.Errorf("vars is not set for env.options")
	}

	if len(config.Options) == 0 {
		return nil, fmt.Errorf("options is not set for env.options")
	}

	options := make(map[string]bool)

	for _, o := range config.Options {
		options[o] = true
	}

	return &EnvOption{config.Vars, options}, nil
}

func (r *EnvOption) Examinate() (bool, error) {
	unset := []string{}
	invalid := []string{}

	for _, v := range r.Vars {
		if val, ok := os.LookupEnv(v); !ok {
			unset = append(unset, v)
		} else if _, ok := r.Options[val]; !ok {
			invalid = append(invalid, v)
		}
	}

	err := ""

	if len(unset) > 0 {
		err += fmt.Sprintf("environment variables not set %v\n", unset)
	}

	if len(invalid) > 0 {
		err += fmt.Sprintf("environment variables not matching options %v\n", invalid)
	}

	if err != "" {
		return false, fmt.Errorf("%v", err)
	}

	return true, nil
}

// Check if an environment variable is a number
// This rule will check if the environment variable is a number
//
// type: env.int,
// vars: []string
type EnvInteger struct {
	Vars []string
}

func (r *EnvInteger) Type() string {
	return "env.int"
}

func (r *EnvInteger) Parse(config config.Exam) (Exam, error) {
	if config.Type != r.Type() {
		return nil, fmt.Errorf("invalid type %v for env.int", config.Type)
	}

	if len(config.Vars) == 0 {
		return nil, fmt.Errorf("vars is not set for env.int")
	}

	return &EnvInteger{config.Vars}, nil
}

func (r *EnvInteger) Examinate() (bool, error) {
	unset := []string{}
	invalid := []string{}

	for _, v := range r.Vars {
		if val, ok := os.LookupEnv(v); !ok {
			unset = append(unset, v)
		} else if _, err := strconv.Atoi(val); err != nil {
			invalid = append(invalid, v)
		}
	}

	err := ""

	if len(unset) > 0 {
		err += fmt.Sprintf("environment variables not set %v\n", unset)
	}

	if len(invalid) > 0 {
		err += fmt.Sprintf("environment variables not set to integer numbers %v\n", invalid)
	}

	if err != "" {
		return false, fmt.Errorf("%v", err)
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
	Vars []string
	Min  int
	Max  int
}

func (r *EnvIntegerRange) Type() string {
	return "env.int-range"
}

func (r *EnvIntegerRange) Parse(config config.Exam) (Exam, error) {
	if config.Type != r.Type() {
		return nil, fmt.Errorf("invalid type %v for env.int-range", config.Type)
	}

	if len(config.Vars) == 0 {
		return nil, fmt.Errorf("vars is not set for env.int-range")
	}

	switch config.Min.(type) {
	case int:
	default:
		return nil, fmt.Errorf("min is not an integer for env.int-range")
	}

	switch config.Max.(type) {
	case int:
	default:
		return nil, fmt.Errorf("max is not an integer for env.int-range")
	}

	return &EnvIntegerRange{config.Vars, config.Min.(int), config.Max.(int)}, nil
}

func (r *EnvIntegerRange) Examinate() (bool, error) {
	unset := []string{}
	not_integer := []string{}
	out_of_bound := []string{}

	for _, v := range r.Vars {
		if val, ok := os.LookupEnv(v); !ok {
			unset = append(unset, v)
		} else if num, err := strconv.Atoi(val); err != nil {
			not_integer = append(not_integer, v)
		} else if num < r.Min || num > r.Max {
			out_of_bound = append(out_of_bound, v)
		}
	}

	err := ""

	if len(unset) > 0 {
		err += fmt.Sprintf("environment variables not set %v\n", unset)
	}

	if len(not_integer) > 0 {
		err += fmt.Sprintf("environment variables not set to integer numbers %v\n", not_integer)
	}

	if len(out_of_bound) > 0 {
		err += fmt.Sprintf("environment variables not in the integer range [%v,%v] %v\n", r.Min, r.Max, out_of_bound)
	}

	if err != "" {
		return false, fmt.Errorf("%v", err)
	}

	return true, nil
}

// Check if an environment variable is a floating number
//
// type: env.float,
// vars: []string
type EnvFloat struct {
	Vars []string
}

func (r *EnvFloat) Type() string {
	return "env.float"
}

func (r *EnvFloat) Parse(config config.Exam) (Exam, error) {
	if config.Type != r.Type() {
		return nil, fmt.Errorf("invalid type %v for env.float", config.Type)
	}

	if len(config.Vars) == 0 {
		return nil, fmt.Errorf("vars is not set for env.float")
	}

	return &EnvFloat{config.Vars}, nil
}

func (r *EnvFloat) Examinate() (bool, error) {
	unset := []string{}
	not_float := []string{}

	for _, v := range r.Vars {
		if val, ok := os.LookupEnv(v); !ok {
			unset = append(unset, v)
		} else if _, err := strconv.ParseFloat(val, 64); err != nil {
			not_float = append(not_float, v)
		}
	}

	err := ""

	if len(unset) > 0 {
		err += fmt.Sprintf("environment variables not set %v\n", unset)
	}

	if len(not_float) > 0 {
		err += fmt.Sprintf("environment variables not set to float numbers %v\n", not_float)
	}

	if err != "" {
		return false, fmt.Errorf("%v", err)
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
	Vars []string
	Min  float64
	Max  float64
}

func (r *EnvFloatRange) Type() string {
	return "env.float-range"
}

func (r *EnvFloatRange) Parse(config config.Exam) (Exam, error) {
	if config.Type != r.Type() {
		return nil, fmt.Errorf("invalid type %v for env.float-range", config.Type)
	}

	if len(config.Vars) == 0 {
		return nil, fmt.Errorf("vars is not set for env.float-range")
	}

	switch config.Min.(type) {
	case float64:
	default:
		return nil, fmt.Errorf("min is not a float for env.float-range")
	}

	switch config.Max.(type) {
	case float64:
	default:
		return nil, fmt.Errorf("max is not a float for env.float-range")
	}

	return &EnvFloatRange{config.Vars, config.Min.(float64), config.Max.(float64)}, nil
}

func (r *EnvFloatRange) Examinate() (bool, error) {
	unset := []string{}
	out_of_bound := []string{}
	not_float := []string{}

	for _, v := range r.Vars {
		val, ok := os.LookupEnv(v)

		if !ok {
			unset = append(unset, v)
		} else if num, err := strconv.ParseFloat(val, 64); err != nil {
			not_float = append(not_float, v)
		} else if num < r.Min || num > r.Max {
			out_of_bound = append(out_of_bound, v)
		}
	}

	err := ""

	if len(unset) > 0 {
		err += fmt.Sprintf("environment variables not set %v\n", unset)
	}

	if len(not_float) > 0 {
		err += fmt.Sprintf("environment variables not set to float numbers %v\n", not_float)
	}

	if len(out_of_bound) > 0 {
		err += fmt.Sprintf("environment variables not in the float range [%v,%v] %v\n", r.Min, r.Max, out_of_bound)
	}

	if err != "" {
		return false, fmt.Errorf("%v", err)
	}

	return true, nil
}

// Check if an environment variable is set to a file that exists
//
// type: env.file,
// vars: []string
type EnvFile struct {
	Vars []string
}

func (r *EnvFile) Type() string {
	return "env.file"
}

func (r *EnvFile) Parse(config config.Exam) (Exam, error) {
	if config.Type != r.Type() {
		return nil, fmt.Errorf("invalid type %v for env.file", config.Type)
	}

	if len(config.Vars) == 0 {
		return nil, fmt.Errorf("vars is not set for env.file")
	}

	return &EnvFile{config.Vars}, nil
}

func (r *EnvFile) Examinate() (bool, error) {
	unset := []string{}
	invalid := []string{}

	for _, v := range r.Vars {
		val, ok := os.LookupEnv(v)
		if !ok {
			unset = append(unset, v)
		} else if _, err := os.Stat(val); err != nil {
			invalid = append(invalid, v)
		}
	}

	err := ""

	if len(unset) > 0 {
		err += fmt.Sprintf("environment variables not set %v\n", unset)
	}

	if len(invalid) > 0 {
		err += fmt.Sprintf("environment variables not pointing to valid files %v\n", invalid)
	}

	if err != "" {
		return false, fmt.Errorf("%v", err)
	}

	return true, nil
}

// Check if an environment variable is set to a directory that exists
//
// type: env.dir,
// vars: []string
type EnvDir struct {
	Vars []string
}

func (r *EnvDir) Type() string {
	return "env.dir"
}

func (r *EnvDir) Parse(config config.Exam) (Exam, error) {
	if config.Type != r.Type() {
		return nil, fmt.Errorf("invalid type %v for env.dir", config.Type)
	}

	if len(config.Vars) == 0 {
		return nil, fmt.Errorf("vars is not set for env.dir")
	}

	return &EnvDir{config.Vars}, nil
}

func (r *EnvDir) Examinate() (bool, error) {
	unset := []string{}
	invalid := []string{}

	for _, v := range r.Vars {
		val, ok := os.LookupEnv(v)
		if !ok {
			unset = append(unset, v)
		} else if stat, err := os.Stat(val); err != nil || !stat.IsDir() {
			invalid = append(invalid, v)
		}
	}

	err := ""

	if len(unset) > 0 {
		err += fmt.Sprintf("environment variables not set %v\n", unset)
	}

	if len(invalid) > 0 {
		err += fmt.Sprintf("environment variables not pointing to valid directories %v\n", invalid)
	}

	if err != "" {
		return false, fmt.Errorf("%v", err)
	}

	return true, nil
}

// Check if an environment variable is set to an IPv4 address
//
// type: env.ipv4,
// vars: []string
type EnvIpv4Addr struct {
	Vars []string
}

func (r *EnvIpv4Addr) Type() string {
	return "env.ipv4"
}

func (r *EnvIpv4Addr) Parse(config config.Exam) (Exam, error) {
	if config.Type != r.Type() {
		return nil, fmt.Errorf("invalid type %v for env.ipv4", config.Type)
	}

	if len(config.Vars) == 0 {
		return nil, fmt.Errorf("vars is not set for env.ipv4")
	}

	return &EnvIpv4Addr{config.Vars}, nil
}

func (r *EnvIpv4Addr) Examinate() (bool, error) {
	unset := []string{}
	invalid := []string{}
	regexp := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)

	for _, v := range r.Vars {
		val, ok := os.LookupEnv(v)
		if !ok {
			unset = append(unset, v)
		} else if !regexp.MatchString(val) {
			invalid = append(invalid, v)
		}
	}

	err := ""

	if len(unset) > 0 {
		err += fmt.Sprintf("environment variables not set %v\n", unset)
	}

	if len(invalid) > 0 {
		err += fmt.Sprintf("environment variables not set to valid IPv4 addresses %v\n", invalid)
	}

	if err != "" {
		return false, fmt.Errorf("%v", err)
	}

	return true, nil
}

// Check if an environment variable is set to an IPv6 address
//
// type: env.ipv6,
// vars: []string
type EnvIpv6Addr struct {
	Vars []string
}

func (r *EnvIpv6Addr) Type() string {
	return "env.ipv6"
}

func (r *EnvIpv6Addr) Parse(config config.Exam) (Exam, error) {
	if config.Type != r.Type() {
		return nil, fmt.Errorf("invalid type %v for env.ipv6", config.Type)
	}

	if len(config.Vars) == 0 {
		return nil, fmt.Errorf("vars is not set for env.ipv6")
	}

	return &EnvIpv6Addr{config.Vars}, nil
}

func (r *EnvIpv6Addr) Examinate() (bool, error) {
	unset := []string{}
	invalid := []string{}
	regexp := regexp.MustCompile(`^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`)

	for _, v := range r.Vars {
		val, ok := os.LookupEnv(v)
		if !ok {
			unset = append(unset, v)
		} else if !regexp.MatchString(val) {
			invalid = append(invalid, v)
		}
	}

	err := ""

	if len(unset) > 0 {
		err += fmt.Sprintf("environment variables not set %v\n", unset)
	}

	if len(invalid) > 0 {
		err += fmt.Sprintf("environment variables not set to valid IPv6 addresses %v\n", invalid)
	}

	if err != "" {
		return false, fmt.Errorf("%v", err)
	}

	return true, nil
}

// Check if an environment variable is set to an IP address
//
// type: env.ip,
// vars: []string
type EnvIpAddr struct {
	Vars []string
}

func (r *EnvIpAddr) Type() string {
	return "env.ip"
}

func (r *EnvIpAddr) Parse(config config.Exam) (Exam, error) {
	if config.Type != r.Type() {
		return nil, fmt.Errorf("invalid type %v for env.ip", config.Type)
	}

	if len(config.Vars) == 0 {
		return nil, fmt.Errorf("vars is not set for env.ip")
	}

	return &EnvIpAddr{config.Vars}, nil
}

// TODO: Refactor this
func (r *EnvIpAddr) Examinate() (bool, error) {
	unset := []string{}
	invalid := []string{}

	for _, v := range r.Vars {
		ipv4 := &EnvIpv4Addr{Vars: []string{v}}
		ipv6 := &EnvIpv6Addr{Vars: []string{v}}

		if ok, _ := ipv4.Examinate(); !ok {
			if ok, _ := ipv6.Examinate(); !ok {
				invalid = append(invalid, v)
			}
		}
	}

	err := ""

	if len(unset) > 0 {
		err += fmt.Sprintf("environment variables not set %v\n", unset)
	}

	if len(invalid) > 0 {
		err += fmt.Sprintf("environment variables not set to valid IP addresses %v\n", invalid)
	}

	if err != "" {
		return false, fmt.Errorf("%v", err)
	}

	return true, nil
}

// A rule to check if an environment variable is set to a hostname
//
// type: env.hostname,
// vars: []string,
// protocol: string
type EnvHostname struct {
	Vars     []string
	Protocol string
}

func (r *EnvHostname) Type() string {
	return "env.hostname"
}

func (r *EnvHostname) Parse(config config.Exam) (Exam, error) {
	if config.Type != r.Type() {
		return nil, fmt.Errorf("invalid type %v for env.hostname", config.Type)
	}

	if len(config.Vars) == 0 {
		return nil, fmt.Errorf("vars is not set for env.hostname")
	}

	return &EnvHostname{config.Vars, config.Protocol}, nil
}

func (r *EnvHostname) Examinate() (bool, error) {
	unset := []string{}
	invalid := []string{}

	for _, v := range r.Vars {
		val, ok := os.LookupEnv(v)
		if !ok {
			unset = append(unset, v)
		} else if ok, _ := r.validateUrl(val); !ok {
			invalid = append(invalid, v)
		}
	}

	err := ""

	if len(unset) > 0 {
		err += fmt.Sprintf("environment variables not set %v\n", unset)
	}

	if len(invalid) > 0 {
		err += fmt.Sprintf("environment variables not set to valid hostnames %v\n", invalid)
	}

	if err != "" {
		return false, fmt.Errorf("%v", err)
	}

	return true, nil
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
		return false, fmt.Errorf("URL %v does not match protocol %v", rawUrl, r.Protocol)
	}

	return true, nil
}
