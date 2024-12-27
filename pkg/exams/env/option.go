package env

import (
	"fmt"

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
	return DefaultExaminate(r.Vars, func(name, value string) (bool, error) {
		if _, ok := r.Options[value]; !ok {
			return false, &InvalidEnvVarError{Var: name, Value: value, Message: r.ErrorMessage()}
		}

		return true, nil
	})
}

func (r *Option) ErrorMessage() string {
	return fmt.Sprintf("value should be one of %v", r.Options)
}
