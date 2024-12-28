package env

import (
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

func (r *NotEmpty) Examinate() exams.Report {
	return DefaultExaminate(r.Vars, func(name, value string) EnvStatus {
		if strings.TrimSpace(value) == "" {
			return invalidEnvVarStatus(name, value, "value must contain at least one non-whitespace character")
		}

		return validEnvVarStatus(name)
	})
}
