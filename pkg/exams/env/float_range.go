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
		return nil, fmt.Errorf("vars is not set for env.float-range")
	}

	switch config.Min.(type) {
	case float64:
	default:
		return nil, fmt.Errorf("min is not a float for env.float-range")
	}

	switch config.Max.(type) {
	case float64:
	default:
		return nil, fmt.Errorf("max is not a float for env.float-range")
	}

	return &FloatRange{config.Vars, config.Min.(float64), config.Max.(float64)}, nil
}

func (r *FloatRange) Examinate() (bool, error) {
	unset := []string{}
	out_of_bound := []string{}
	not_float := []string{}

	for _, v := range r.Vars {
		val, ok := os.LookupEnv(v)

		if !ok {
			unset = append(unset, v)
		} else if num, err := strconv.ParseFloat(val, 64); err != nil {
			not_float = append(not_float, v)
		} else if num < r.Min || num > r.Max {
			out_of_bound = append(out_of_bound, v)
		}
	}

	err := ""

	if len(unset) > 0 {
		err += fmt.Sprintf("environment variables not set %v\n", unset)
	}

	if len(not_float) > 0 {
		err += fmt.Sprintf("environment variables not set to float numbers %v\n", not_float)
	}

	if len(out_of_bound) > 0 {
		err += fmt.Sprintf("environment variables not in the float range [%v,%v] %v\n", r.Min, r.Max, out_of_bound)
	}

	if err != "" {
		return false, fmt.Errorf("%v", err)
	}

	return true, nil
}
