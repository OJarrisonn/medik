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

	switch config.Min.(type) {
	case int:
	default:
		return nil, fmt.Errorf("min is not an integer for env.int-range")
	}

	switch config.Max.(type) {
	case int:
	default:
		return nil, fmt.Errorf("max is not an integer for env.int-range")
	}

	return &IntRange{config.Vars, config.Min.(int), config.Max.(int)}, nil
}

func (r *IntRange) Examinate() (bool, error) {
	unset := []string{}
	not_integer := []string{}
	out_of_bound := []string{}

	for _, v := range r.Vars {
		if val, ok := os.LookupEnv(v); !ok {
			unset = append(unset, v)
		} else if num, err := strconv.Atoi(val); err != nil {
			not_integer = append(not_integer, v)
		} else if num < r.Min || num > r.Max {
			out_of_bound = append(out_of_bound, v)
		}
	}

	err := ""

	if len(unset) > 0 {
		err += fmt.Sprintf("environment variables not set %v\n", unset)
	}

	if len(not_integer) > 0 {
		err += fmt.Sprintf("environment variables not set to integer numbers %v\n", not_integer)
	}

	if len(out_of_bound) > 0 {
		err += fmt.Sprintf("environment variables not in the integer range [%v,%v] %v\n", r.Min, r.Max, out_of_bound)
	}

	if err != "" {
		return false, fmt.Errorf("%v", err)
	}

	return true, nil
}
