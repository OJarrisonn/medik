package env

import (
	"os"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
)

// Check if a given environment variable is set
//
// type: env.is-set,
// vars: []string
type IsSet struct {
	Vars []string
}

func (r *IsSet) Type() string {
	return "env.is-set"
}

func (r *IsSet) Parse(config config.Exam) (exams.Exam, error) {
	if config.Type != r.Type() {
		return nil, &exams.WrongExamParserError{Source: config.Type, Using: r.Type()}
	}

	if len(config.Vars) == 0 {
		return nil, &VarsUnsetError{Exam: r.Type()}
	}

	return &IsSet{config.Vars}, nil
}

func (r *IsSet) Examinate() (bool, []error) {
	errors := make([]error, len(r.Vars))
	hasError := false

	for i, v := range r.Vars {
		if _, ok := os.LookupEnv(v); !ok {
			hasError = true
			errors[i] = &UnsetEnvVarError{Var: v}
		}
	}

	if hasError {
		return false, errors
	}

	return true, nil
}
