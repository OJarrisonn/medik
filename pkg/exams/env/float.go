package env

import (
	"fmt"
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
		return nil, &exams.WrongExamParserError{Got: config.Type, Expected: r.Type()}
	}

	if len(config.Vars) == 0 {
		return nil, fmt.Errorf("vars is not set for env.float")
	}

	return &Float{config.Vars}, nil
}

func (r *Float) Examinate() (bool, error) {
	unset := []string{}
	not_float := []string{}

	for _, v := range r.Vars {
		if val, ok := os.LookupEnv(v); !ok {
			unset = append(unset, v)
		} else if _, err := strconv.ParseFloat(val, 64); err != nil {
			not_float = append(not_float, v)
		}
	}

	err := ""

	if len(unset) > 0 {
		err += fmt.Sprintf("environment variables not set %v\n", unset)
	}

	if len(not_float) > 0 {
		err += fmt.Sprintf("environment variables not set to float numbers %v\n", not_float)
	}

	if err != "" {
		return false, fmt.Errorf("%v", err)
	}

	return true, nil
}
