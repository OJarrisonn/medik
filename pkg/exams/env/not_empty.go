package env

import (
	"os"
	"strings"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
)

// Check if an environment variable is set and not empty
// Strings with only whitespace are considered empty
//
// type: env.not-empty
// vars: []string
type NotEmpty struct {
	Vars []string
}

func (r *NotEmpty) Type() string {
	return "env.not-empty"
}

func (r *NotEmpty) Parse(config config.Exam) (exams.Exam, error) {
	if config.Type != r.Type() {
		return nil, &exams.WrongExamParserError{Source: config.Type, Using: r.Type()}
	}

	if len(config.Vars) == 0 {
		return nil, &VarsUnsetError{Exam: r.Type()}
	}

	return &NotEmpty{config.Vars}, nil
}

func (r *NotEmpty) Examinate() (bool, []error) {
	errors := make([]error, len(r.Vars))
	hasError := false

	for i, v := range r.Vars {
		if val, ok := os.LookupEnv(v); !ok {
			hasError = true
			errors[i] = &UnsetEnvVarError{Var: v}
		} else if strings.TrimSpace(val) == "" {
			hasError = true
			errors[i] = &InvalidEnvVarError{Var: v, Value: val, Message: "value must contain at least one non-whitespace character"}
		}
	}

	if hasError {
		return false, errors
	}

	return true, nil
}
