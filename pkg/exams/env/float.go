package env

import (
	"strconv"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
	"github.com/OJarrisonn/medik/pkg/medik"
)

// Check if an environment variable is a floating number
//
// type: env.float,
// vars: []string
type Float struct {
	Vars  []string
	Level int
}

func (r *Float) Type() string {
	return "env.float"
}

func (r *Float) Parse(conf config.Exam) (exams.Exam, error) {
	return DefaultParse[*Float](conf, func(conf config.Exam) (exams.Exam, error) {
		return &Float{conf.Vars, medik.LogLevelFromStr(conf.Level)}, nil
	})
}

func (r *Float) Examinate() exams.Report {
	return DefaultExaminate(r.Type(), r.Level, r.Vars, func(name, value string) EnvStatus {
		_, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return invalidEnvVarStatus(name, r.Level, value, err.Error())
		}

		return validEnvVarStatus(name)
	})
}
