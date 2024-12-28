package env

import (
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

func (r *Int) Examinate() exams.Report {
	return DefaultExaminate(r.Vars, func(name, value string) EnvStatus {
		_, err := strconv.Atoi(value)

		if err != nil {
			return invalidEnvVarStatus(name, value, r.ErrorMessage(err))
		}

		return validEnvVarStatus(name)
	})
}

func (r *Int) ErrorMessage(err error) string {
	return "value should be a number. " + err.Error()
}
