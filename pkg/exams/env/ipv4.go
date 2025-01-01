package env

import (
	"regexp"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
	"github.com/OJarrisonn/medik/pkg/medik"
)

// Check if an environment variable is set to an IPv4 address
//
// type: env.ipv4,
// vars: []string
type Ipv4 struct {
	Vars  []string
	Level int
}

func (r *Ipv4) Type() string {
	return "env.ipv4"
}

func (r *Ipv4) Parse(conf config.Exam) (exams.Exam, error) {
	return DefaultParse[*Ipv4](conf, func(config config.Exam) (exams.Exam, error) {
		return &Ipv4{config.Vars, medik.LogLevelFromStr(config.Level)}, nil
	})
}

func (r *Ipv4) Examinate() exams.Report {
	return DefaultExaminate(r.Type(), r.Level, r.Vars, func(name, value string) EnvStatus {
		regexp := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)

		if !regexp.MatchString(value) {
			return invalidEnvVarStatus(name, r.Level, value, "value should be a valid IPv4 address")
		}

		return validEnvVarStatus(name)
	})
}
