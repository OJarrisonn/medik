package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
)

// Check if an environment variable is set and not empty
// Strings with only whitespace are considered empty
//
// type: env.not-empty
// vars: []string
type NotEmpty struct {
	Vars []string
}

func (r *NotEmpty) Type() string {
	return "env.not-empty"
}

func (r *NotEmpty) Parse(config config.Exam) (exams.Exam, error) {
	if config.Type != r.Type() {
		return nil, &exams.WrongExamParserError{Source: config.Type, Using: r.Type()}
	}

	if len(config.Vars) == 0 {
		return nil, &VarsUnsetError{Exam: r.Type()}
	}

	return &NotEmpty{config.Vars}, nil
}

func (r *NotEmpty) Examinate() (bool, error) {
	unset := []string{}
	invalid := []string{}

	for _, v := range r.Vars {
		if val, ok := os.LookupEnv(v); !ok {
			unset = append(unset, v)
		} else if strings.TrimSpace(val) == "" {
			invalid = append(invalid, v)
		}
	}

	err := ""

	if len(unset) > 0 {
		err += fmt.Sprintf("environment variables not set %v\n", unset)
	}

	if len(invalid) > 0 {
		err += fmt.Sprintf("environment variables set to empty strings %v\n", invalid)
	}

	if err != "" {
		return false, fmt.Errorf("%v", err)
	}

	return true, nil
}
