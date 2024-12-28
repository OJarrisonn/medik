package env

import (
	"fmt"
	"os"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
)

// Check if an environment variable is set to a file that exists
//
// type: env.file,
// vars: []string
type File struct {
	Vars   []string
	Exists bool
}

func (r *File) Type() string {
	return "env.file"
}

func (r *File) Parse(config config.Exam) (exams.Exam, error) {
	if config.Type != r.Type() {
		return nil, &exams.WrongExamParserError{Source: config.Type, Using: r.Type()}
	}

	if len(config.Vars) == 0 {
		return nil, &VarsUnsetError{Exam: r.Type()}
	}

	return &File{config.Vars, config.Exists}, nil
}

func (r *File) Examinate() (bool, []error) {
	return DefaultExaminate(r.Vars, func(name, value string) (bool, error) {
		_, err := os.Stat(value)

		if (err == nil) != r.Exists {
			return false, &InvalidEnvVarError{Var: name, Value: value, Message: r.ErrorMessage(err)}
		}

		return true, nil
	})
}

func (r *File) ErrorMessage(err error) string {
	non := "an "
	if !r.Exists {
		non = "a non "
	}

	return fmt.Sprintf("value should point to %vexisting file. %v", non, err)
}
