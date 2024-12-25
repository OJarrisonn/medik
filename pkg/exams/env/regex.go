package env

import (
	"fmt"
	"os"
	"regexp"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
)

// Check if an environment variable is set and matches a regex
// The regex is a string that will be compiled into a regular expression using Go's regexp package
//
// type: env.regex,
// vars: []string,
// regex: string
type Regex struct {
	Vars  []string
	Regex *regexp.Regexp
}

func (r *Regex) Type() string {
	return "env.regex"
}

func (r *Regex) Parse(config config.Exam) (exams.Exam, error) {
	if config.Type != r.Type() {
		return nil, &exams.WrongExamParserError{Source: config.Type, Using: r.Type()}
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

	return &Regex{config.Vars, regexp}, nil
}

func (r *Regex) Examinate() (bool, error) {
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
