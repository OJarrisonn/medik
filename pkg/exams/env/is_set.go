package env

import (
	"fmt"
	"os"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
)

// Check if a given environment variable is set
//
// type: env.is-set,
// vars: []string
type IsSet struct {
	Vars []string
}

func (r *IsSet) Type() string {
	return "env.is-set"
}

func (r *IsSet) Parse(config config.Exam) (exams.Exam, error) {
	if config.Type != r.Type() {
		return nil, &exams.WrongExamParserError{Got: config.Type, Expected: r.Type()}
	}

	if len(config.Vars) == 0 {
		return nil, fmt.Errorf("vars is not set for env.is-set")
	}

	return &IsSet{config.Vars}, nil
}

func (r *IsSet) Examinate() (bool, error) {
	unset := []string{}

	for _, v := range r.Vars {
		if _, ok := os.LookupEnv(v); !ok {
			unset = append(unset, v)
		}
	}

	if len(unset) > 0 {
		return false, fmt.Errorf("environment variables not set %v", unset)
	}

	return true, nil
}
