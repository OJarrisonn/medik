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

func (r *IntRange) Examinate() exams.Report {
	return DefaultExaminate(r.Type(), r.Vars, func(name, value string) EnvStatus {
		num, err := strconv.Atoi(value)

		if err != nil {
			return invalidEnvVarStatus(name, value, "value should be a number. "+err.Error())
		}

		if num < r.Min || num > r.Max {
			return invalidEnvVarStatus(name, value, fmt.Sprintf("value should be in the range [%v,%v]", r.Min, r.Max))
		}

		return validEnvVarStatus(name)
	})
}
