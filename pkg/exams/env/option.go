package env

import (
	"fmt"
	"os"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
)

// Check if an environment variable is set and matches one of a list of possible values
//
// type: env.options,
// vars: []string,
// options: []string
type Option struct {
	Vars    []string
	Options map[string]bool
}

func (r *Option) Type() string {
	return "env.options"
}

func (r *Option) Parse(config config.Exam) (exams.Exam, error) {
	if config.Type != r.Type() {
		return nil, &exams.WrongExamParserError{Source: config.Type, Using: r.Type()}
	}

	if len(config.Vars) == 0 {
		return nil, &VarsUnsetError{Exam: r.Type()}
	}

	if len(config.Options) == 0 {
		return nil, fmt.Errorf("options is not set for env.options")
	}

	options := make(map[string]bool)

	for _, o := range config.Options {
		options[o] = true
	}

	return &Option{config.Vars, options}, nil
}

func (r *Option) Examinate() (bool, error) {
	unset := []string{}
	invalid := []string{}

	for _, v := range r.Vars {
		if val, ok := os.LookupEnv(v); !ok {
			unset = append(unset, v)
		} else if _, ok := r.Options[val]; !ok {
			invalid = append(invalid, v)
		}
	}

	err := ""

	if len(unset) > 0 {
		err += fmt.Sprintf("environment variables not set %v\n", unset)
	}

	if len(invalid) > 0 {
		err += fmt.Sprintf("environment variables not matching options %v\n", invalid)
	}

	if err != "" {
		return false, fmt.Errorf("%v", err)
	}

	return true, nil
}
