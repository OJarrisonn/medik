package env

import (
	"fmt"
	"os"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
	"github.com/OJarrisonn/medik/pkg/medik"
)

// Check if an environment variable is set to a directory that exists
//
// type: env.dir,
// vars: []string
type Dir struct {
	Vars   []string
	Level  int
	Exists bool
}

func (r *Dir) Type() string {
	return "env.dir"
}

func (r *Dir) Parse(conf config.Exam) (exams.Exam, error) {
	return DefaultParse[*Dir](conf, func(conf config.Exam) (exams.Exam, error) {
		return &Dir{conf.Vars, medik.LogLevelFromStr(conf.Level), conf.Exists}, nil
	})
}

func (r *Dir) Examinate() exams.Report {
	return DefaultExaminate(r.Type(), r.Level, r.Vars, func(name, value string) EnvStatus {
		stat, err := os.Stat(value)

		if exists := err == nil && stat.IsDir(); exists != r.Exists {
			return invalidEnvVarStatus(name, r.Level, value, r.ErrorMessage(err))
		}

		return validEnvVarStatus(name)
	})
}

func (r *Dir) ErrorMessage(err error) string {
	non := "an "
	if !r.Exists {
		non = "a non "
	}

	return fmt.Sprintf("value should point to %vexisting directory. %v", non, err)
}
