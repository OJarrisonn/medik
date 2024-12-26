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
		return nil, &VarsUnsetError{Exam: r.Type()}
	}

	if config.Regex == "" {
		return nil, &exams.MissingFieldError{Field: "regex", Exam: r.Type()}
	}

	regexp, rerr := regexp.Compile(config.Regex)

	if rerr != nil {
		return nil, &exams.FieldValueError{Field: "regex", Exam: r.Type(), Value: config.Regex, Message: rerr.Error()}
	}

	return &Regex{config.Vars, regexp}, nil
}

func (r *Regex) Examinate() (bool, []error) {
	errors := make([]error, len(r.Vars))
	hasError := false

	for i, v := range r.Vars {
		if val, ok := os.LookupEnv(v); !ok {
			hasError = true
			errors[i] = &UnsetEnvVarError{Var: v}
		} else if !r.Regex.MatchString(val) {
			hasError = true
			errors[i] = &InvalidEnvVarError{Var: v, Value: val, Message: r.ErrorMessage()}
		}
	}

	if hasError {
		return false, errors
	}

	return true, nil
}

func (r *Regex) ErrorMessage() string {
	return fmt.Sprintf("value should match regex %v", r.Regex)
}