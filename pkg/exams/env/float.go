package env

import (
	"os"
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
	errors := make([]error, len(r.Vars))
	hasError := false

	for i, v := range r.Vars {
		if val, ok := os.LookupEnv(v); !ok {
			hasError = true
			errors[i] = &UnsetEnvVarError{Var: v}
		} else if _, err := strconv.ParseFloat(val, 64); err != nil {
			hasError = true
			errors[i] = &InvalidEnvVarError{Var: v, Value: val, Message: err.Error()}
		}
	}

	if hasError {
		return false, errors
	}

	return true, nil
}
