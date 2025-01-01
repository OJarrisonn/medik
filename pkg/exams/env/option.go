package env

import (
	"fmt"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
	"github.com/OJarrisonn/medik/pkg/medik"
)

// Check if an environment variable is set and matches one of a list of possible values
//
// type: env.options,
// vars: []string,
// options: []string
type Option struct {
	Vars    []string
	Level   int
	Options map[string]bool
}

func (r *Option) Type() string {
	return "env.options"
}

func (r *Option) Parse(conf config.Exam) (exams.Exam, error) {
	return DefaultParse[*Option](conf, func(config config.Exam) (exams.Exam, error) {
		if len(config.Options) == 0 {
			return nil, &exams.MissingFieldError{Field: "options", Exam: r.Type()}
		}

		options := make(map[string]bool)

		for _, o := range config.Options {
			options[o] = true
		}

		return &Option{config.Vars, medik.LogLevelFromStr(config.Level), options}, nil
	})
}

func (r *Option) Examinate() exams.Report {
	return DefaultExaminate(r.Type(), r.Level, r.Vars, func(name, value string) EnvStatus {
		if _, ok := r.Options[value]; !ok {
			return invalidEnvVarStatus(name, r.Level, value, r.ErrorMessage())
		}

		return validEnvVarStatus(name)
	})
}

func (r *Option) ErrorMessage() string {
	return fmt.Sprintf("value should be one of %v", r.Options)
}
