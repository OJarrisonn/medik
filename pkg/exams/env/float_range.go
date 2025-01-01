package env

import (
	"fmt"
	"strconv"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
	"github.com/OJarrisonn/medik/pkg/medik"
)

// Check if an environment variable is a number within a range
// Min and Max are inclusive
//
// type: env.float-range,
// vars: []string,
// min: float,
// max: float
type FloatRange struct {
	Vars  []string
	Level int
	Min   float64
	Max   float64
}

func (r *FloatRange) Type() string {
	return "env.float-range"
}

func (r *FloatRange) Parse(conf config.Exam) (exams.Exam, error) {
	return DefaultParse[*FloatRange](conf, func(conf config.Exam) (exams.Exam, error) {
		if conf.Min == nil {
			return nil, &exams.MissingFieldError{Field: "min", Exam: r.Type()}
		}

		switch conf.Min.(type) {
		case float64:
		default:
			return nil, &exams.FieldValueError{Field: "min", Exam: r.Type(), Value: fmt.Sprint(conf.Min), Message: "expected a float value"}
		}

		if conf.Max == nil {
			return nil, &exams.MissingFieldError{Field: "max", Exam: r.Type()}
		}

		switch conf.Max.(type) {
		case float64:
		default:
			return nil, &exams.FieldValueError{Field: "max", Exam: r.Type(), Value: fmt.Sprint(conf.Max), Message: "expected a float value"}
		}

		return &FloatRange{conf.Vars, medik.LogLevelFromStr(conf.Level), conf.Min.(float64), conf.Max.(float64)}, nil
	})
}

func (r *FloatRange) Examinate() exams.Report {
	return DefaultExaminate(r.Type(), r.Level, r.Vars, func(name string, value string) EnvStatus {
		num, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return invalidEnvVarStatus(name, r.Level, value, err.Error())
		}

		if num < r.Min || num > r.Max {
			return invalidEnvVarStatus(name, r.Level, value, fmt.Sprintf("value should be in the range [%v,%v]", r.Min, r.Max))
		}

		return validEnvVarStatus(name)
	})
}
