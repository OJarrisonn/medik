package env

import (
	"fmt"
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

func (r *FloatRange) Examinate() exams.Report {
	return DefaultExaminate(r.Type(), r.Vars, func(name string, value string) EnvStatus {
		num, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return invalidEnvVarStatus(name, value, err.Error())
		}

		if num < r.Min || num > r.Max {
			return invalidEnvVarStatus(name, value, fmt.Sprintf("value should be in the range [%v,%v]", r.Min, r.Max))
		}

		return validEnvVarStatus(name)
	})
}
