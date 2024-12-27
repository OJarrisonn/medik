package env

import (
	"strconv"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
)

// Check if an environment variable is a floating number
//
// type: env.float,
// vars: []string
type Float struct {
	Vars []string
}

func (r *Float) Type() string {
	return "env.float"
}

func (r *Float) Parse(config config.Exam) (exams.Exam, error) {
	if config.Type != r.Type() {
		return nil, &exams.WrongExamParserError{Source: config.Type, Using: r.Type()}
	}

	if len(config.Vars) == 0 {
		return nil, &VarsUnsetError{Exam: r.Type()}
	}

	return &Float{config.Vars}, nil
}

func (r *Float) Examinate() (bool, []error) {
	return DefaultExaminate(r.Vars, func(name, value string) (bool, error) {
		_, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return false, &InvalidEnvVarError{Var: name, Value: value, Message: err.Error()}
		}

		return true, nil
	})
}
