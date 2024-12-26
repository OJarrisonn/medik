package env

import (
	"os"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
)

// Check if an environment variable is set to a directory that exists
//
// type: env.dir,
// vars: []string
type Dir struct {
	Vars []string
}

func (r *Dir) Type() string {
	return "env.dir"
}

func (r *Dir) Parse(config config.Exam) (exams.Exam, error) {
	if config.Type != r.Type() {
		return nil, &exams.WrongExamParserError{Source: config.Type, Using: r.Type()}
	}

	if len(config.Vars) == 0 {
		return nil, &VarsUnsetError{Exam: r.Type()}
	}

	return &Dir{config.Vars}, nil
}

func (r *Dir) Examinate() (bool, []error) {
	errors := make([]error, len(r.Vars))
	hasError := false

	for i, v := range r.Vars {
		val, ok := os.LookupEnv(v)
		if !ok {
			hasError = true
			errors[i] = &UnsetEnvVarError{Var: v}
		} else if stat, err := os.Stat(val); err != nil || !stat.IsDir() {
			hasError = true
			errors[i] = &InvalidEnvVarError{Var: v, Value: val, Message: r.ErrorMessage(err)}
		}
	}

	if hasError {
		return false, errors
	}

	return true, nil
}

func (r *Dir) ErrorMessage(err error) string {
	return "value should point to a existing directory. " + err.Error()
}