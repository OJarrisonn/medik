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
	return DefaultExaminate(r.Vars, func(name, value string) (bool, error) {
		stat, err := os.Stat(value)
		if err != nil || !stat.IsDir() {
			return false, &InvalidEnvVarError{Var: name, Value: value, Message: r.ErrorMessage(err)}
		}

		return true, nil
	})
}

func (r *Dir) ErrorMessage(err error) string {
	return "value should point to a existing directory. " + err.Error()
}
