package env

import (
	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
	"github.com/OJarrisonn/medik/pkg/medik"
)

// Check if a given environment variable is set
//
// type: env.is-set,
// vars: []string
type IsSet struct {
	Vars  []string
	Level int
}

func (r *IsSet) Type() string {
	return "env.is-set"
}

func (r *IsSet) Parse(conf config.Exam) (exams.Exam, error) {
	return DefaultParse[*IsSet](conf, func(config config.Exam) (exams.Exam, error) {
		return &IsSet{config.Vars, medik.LogLevelFromStr(config.Level)}, nil
	})
}

func (r *IsSet) Examinate() exams.Report {
	return DefaultExaminate(r.Type(), r.Level, r.Vars, func(name, value string) EnvStatus {
		return validEnvVarStatus(name)
	})
}
