package env

import (
	"strconv"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
	"github.com/OJarrisonn/medik/pkg/medik"
)

// Check if an environment variable is a number
// This rule will check if the environment variable is a number
//
// type: env.int,
// vars: []string
type Int struct {
	Vars  []string
	Level int
}

func (r *Int) Type() string {
	return "env.int"
}

func (r *Int) Parse(conf config.Exam) (exams.Exam, error) {
	return DefaultParse[*Int](conf, func(conf config.Exam) (exams.Exam, error) {
		return &Int{conf.Vars, medik.LogLevelFromStr(conf.Level)}, nil
	})
}

func (r *Int) Examinate() exams.Report {
	return DefaultExaminate(r.Type(), r.Level, r.Vars, func(name, value string) EnvStatus {
		_, err := strconv.Atoi(value)
		if err != nil {
			return invalidEnvVarStatus(name, r.Level, value, r.ErrorMessage(err))
		}

		return validEnvVarStatus(name)
	})
}

func (r *Int) ErrorMessage(err error) string {
	return "value should be a number. " + err.Error()
}
