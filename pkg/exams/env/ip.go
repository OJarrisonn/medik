package env

import (
	"regexp"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
	"github.com/OJarrisonn/medik/pkg/medik"
)

// Check if an environment variable is set to an IP address
//
// type: env.ip,
// vars: []string
type Ip struct {
	Vars  []string
	Level int
}

func (r *Ip) Type() string {
	return "env.ip"
}

func (r *Ip) Parse(conf config.Exam) (exams.Exam, error) {
	return DefaultParse[*Ip](conf, func(config config.Exam) (exams.Exam, error) {
		return &Ip{config.Vars, medik.LogLevelFromStr(config.Level)}, nil
	})
}

// TODO: Refactor this
func (r *Ip) Examinate() exams.Report {
	return DefaultExaminate(r.Type(), r.Level, r.Vars, func(name, value string) EnvStatus {
		regexpv4 := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)

		regexpv6 := regexp.MustCompile(`^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`)

		if !regexpv4.MatchString(value) && !regexpv6.MatchString(value) {
			return invalidEnvVarStatus(name, r.Level, value, "value should be a valid IP address")
		}

		return validEnvVarStatus(name)
	})
}
