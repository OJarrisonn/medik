package env

import (
	"fmt"
	"os"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
)

// Check if an environment variable is set and matches one of a list of possible values
//
// type: env.options,
// vars: []string,
// options: []string
type Option struct {
	Vars    []string
	Options map[string]bool
}

func (r *Option) Type() string {
	return "env.options"
}

func (r *Option) Parse(config config.Exam) (exams.Exam, error) {
	if config.Type != r.Type() {
		return nil, &exams.WrongExamParserError{Source: config.Type, Using: r.Type()}
	}

	if len(config.Vars) == 0 {
		return nil, &VarsUnsetError{Exam: r.Type()}
	}

	if len(config.Options) == 0 {
		return nil, &exams.MissingFieldError{Field: "options", Exam: r.Type()}
	}

	options := make(map[string]bool)

	for _, o := range config.Options {
		options[o] = true
	}

	return &Option{config.Vars, options}, nil
}

func (r *Option) Examinate() (bool, []error) {
	errors := make([]error, len(r.Vars))
	hasError := false

	for i, v := range r.Vars {
		if val, ok := os.LookupEnv(v); !ok {
			hasError = true
			errors[i] = &UnsetEnvVarError{Var: v}
		} else if _, ok := r.Options[val]; !ok {
			hasError = true
			errors[i] = &InvalidEnvVarError{Var: v, Value: val, Message: r.ErrorMessage()}
		}
	}

	if hasError {
		return false, errors
	}

	return true, nil
}

func (r *Option) ErrorMessage() string {
	return fmt.Sprintf("value should be one of %v", r.Options)
}