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
// type: env.float-range,
// vars: []string,
// min: float,
// max: float
type FloatRange struct {
	Vars []string
	Min  float64
	Max  float64
}

func (r *FloatRange) Type() string {
	return "env.float-range"
}

func (r *FloatRange) Parse(config config.Exam) (exams.Exam, error) {
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
	case float64:
	default:
		return nil, &exams.FieldValueError{Field: "min", Exam: r.Type(), Value: fmt.Sprint(config.Min), Message: "expected a float value"}
	}

	if config.Max == nil {
		return nil, &exams.MissingFieldError{Field: "max", Exam: r.Type()}
	}

	switch config.Max.(type) {
	case float64:
	default:
		return nil, &exams.FieldValueError{Field: "max", Exam: r.Type(), Value: fmt.Sprint(config.Max), Message: "expected a float value"}
	}

	return &FloatRange{config.Vars, config.Min.(float64), config.Max.(float64)}, nil
}

func (r *FloatRange) Examinate() (bool, []error) {
	errors := make([]error, len(r.Vars))
	hasError := false

	for i, v := range r.Vars {
		val, ok := os.LookupEnv(v)

		if !ok {
			hasError = true
			errors[i] = &UnsetEnvVarError{Var: v}
		} else if num, err := strconv.ParseFloat(val, 64); err != nil {
			hasError = true
			errors[i] = &InvalidEnvVarError{Var: v, Value: val, Message: err.Error()}
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
