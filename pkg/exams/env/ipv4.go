package env

import (
	"regexp"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
)

// Check if an environment variable is set to an IPv4 address
//
// type: env.ipv4,
// vars: []string
type Ipv4 struct {
	Vars []string
}

func (r *Ipv4) Type() string {
	return "env.ipv4"
}

func (r *Ipv4) Parse(config config.Exam) (exams.Exam, error) {
	if config.Type != r.Type() {
		return nil, &exams.WrongExamParserError{Source: config.Type, Using: r.Type()}
	}

	if len(config.Vars) == 0 {
		return nil, &VarsUnsetError{Exam: r.Type()}
	}

	return &Ipv4{config.Vars}, nil
}

func (r *Ipv4) Examinate() (bool, []error) {
	return DefaultExaminate(r.Vars, func(name, value string) (bool, error) {
		regexp := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)

		if !regexp.MatchString(value) {
			return false, &InvalidEnvVarError{Var: name, Value: value, Message: "value should be a valid IPv4 address"}
		}

		return true, nil
	})
}
