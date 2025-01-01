package env

import (
	"fmt"
	"os"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
	"github.com/OJarrisonn/medik/pkg/medik"
)

// Check if an environment variable is set to a file that exists
//
// type: env.file,
// vars: []string
type File struct {
	Vars   []string
	Level  int
	Exists bool
}

func (r *File) Type() string {
	return "env.file"
}

func (r *File) Parse(conf config.Exam) (exams.Exam, error) {
	return DefaultParse[*File](conf, func(conf config.Exam) (exams.Exam, error) {
		return &File{conf.Vars, medik.LogLevelFromStr(conf.Level), conf.Exists}, nil
	})
}

func (r *File) Examinate() exams.Report {
	return DefaultExaminate(r.Type(), r.Level, r.Vars, func(name, value string) EnvStatus {
		_, err := os.Stat(value)

		if (err == nil) != r.Exists {
			return invalidEnvVarStatus(name, r.Level, value, r.ErrorMessage(err))
		}

		return validEnvVarStatus(name)
	})
}

func (r *File) ErrorMessage(err error) string {
	non := "an "
	if !r.Exists {
		non = "a non "
	}

	return fmt.Sprintf("value should point to %vexisting file. %v", non, err)
}
