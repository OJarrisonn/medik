package env

import (
	"os"
	"strconv"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
)

// Check if an environment variable is a number
// This rule will check if the environment variable is a number
//
// type: env.int,
// vars: []string
type Int struct {
	Vars []string
}

func (r *Int) Type() string {
	return "env.int"
}

func (r *Int) Parse(config config.Exam) (exams.Exam, error) {
	if config.Type != r.Type() {
		return nil, &exams.WrongExamParserError{Source: config.Type, Using: r.Type()}
	}

	if len(config.Vars) == 0 {
		return nil, &VarsUnsetError{Exam: r.Type()}
	}

	return &Int{config.Vars}, nil
}

func (r *Int) Examinate() (bool, []error) {
	errors := make([]error, len(r.Vars))
	hasError := false

	for i, v := range r.Vars {
		if val, ok := os.LookupEnv(v); !ok {
			hasError = true
			errors[i] = &UnsetEnvVarError{Var: v}
		} else if _, err := strconv.Atoi(val); err != nil {
			hasError = true
			errors[i] = &InvalidEnvVarError{Var: v, Value: val, Message: r.ErrorMessage(err)}
		}
	}

	if hasError {
		return false, errors
	}

	return true, nil
}

func (r *Int) ErrorMessage(err error) string {
	return "value should be a number. " + err.Error()
}
