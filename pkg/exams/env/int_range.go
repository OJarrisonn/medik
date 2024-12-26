package env

import (
	"fmt"
	"os"
	"strconv"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
)

// Check if an environment variable is a number within a range
// Min and Max are inclusive
//
// type: env.int-range,
// vars: []string,
// min: int,
// max: int
type IntRange struct {
	Vars []string
	Min  int
	Max  int
}

func (r *IntRange) Type() string {
	return "env.int-range"
}

func (r *IntRange) Parse(config config.Exam) (exams.Exam, error) {
	if config.Type != r.Type() {
		return nil, &exams.WrongExamParserError{Source: config.Type, Using: r.Type()}
	}

	if len(config.Vars) == 0 {
		return nil, &VarsUnsetError{Exam: r.Type()}
	}

	if config.Min == nil {
		return nil, &exams.MissingFieldError{Field: "min", Exam: r.Type()}
	}

	switch config.Min.(type) {
	case int:
	default:
		return nil, &exams.FieldValueError{Field: "min", Exam: r.Type(), Value: fmt.Sprint(config.Min), Message: "expected an integer value"}
	}

	if config.Max == nil {
		return nil, &exams.MissingFieldError{Field: "max", Exam: r.Type()}
	}

	switch config.Max.(type) {
	case int:
	default:
		return nil, &exams.FieldValueError{Field: "max", Exam: r.Type(), Value: fmt.Sprint(config.Max), Message: "expected an integer value"}
	}

	return &IntRange{config.Vars, config.Min.(int), config.Max.(int)}, nil
}

func (r *IntRange) Examinate() (bool, []error) {
	errors := make([]error, len(r.Vars))
	hasError := false

	for i, v := range r.Vars {
		if val, ok := os.LookupEnv(v); !ok {
			hasError = true
			errors[i] = &UnsetEnvVarError{Var: v}
		} else if num, err := strconv.Atoi(val); err != nil {
			hasError = true
			errors[i] = &InvalidEnvVarError{Var: v, Value: val, Message: r.ErrorMessage(err)}
		} else if num < r.Min || num > r.Max {
			hasError = true
			errors[i] = &InvalidEnvVarError{Var: v, Value: val, Message: fmt.Sprintf("value should be in the range [%v,%v]", r.Min, r.Max)}
		}
	}

	if hasError {
		return false, errors
	}

	return true, nil
}

func (r *IntRange) ErrorMessage(err error) string {
	return "value should be an integer. " + err.Error()
}
