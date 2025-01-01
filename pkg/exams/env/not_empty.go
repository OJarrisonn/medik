package env

import (
	"strings"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
	"github.com/OJarrisonn/medik/pkg/medik"
)

// Check if an environment variable is set and not empty
// Strings with only whitespace are considered empty
//
// type: env.not-empty
// vars: []string
type NotEmpty struct {
	Vars  []string
	Level int
}

func (r *NotEmpty) Type() string {
	return "env.not-empty"
}

func (r *NotEmpty) Parse(conf config.Exam) (exams.Exam, error) {
	return DefaultParse[*NotEmpty](conf, func(config config.Exam) (exams.Exam, error) {
		return &NotEmpty{config.Vars, medik.LogLevelFromStr(config.Level)}, nil
	})
}

func (r *NotEmpty) Examinate() exams.Report {
	return DefaultExaminate(r.Type(), r.Level, r.Vars, func(name, value string) EnvStatus {
		if strings.TrimSpace(value) == "" {
			return invalidEnvVarStatus(name, r.Level, value, "value must contain at least one non-whitespace character")
		}

		return validEnvVarStatus(name)
	})
}
