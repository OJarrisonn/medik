package env

import (
	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
)

// Check if an environment variable is set to an IP address
//
// type: env.ip,
// vars: []string
type Ip struct {
	Vars []string
}

func (r *Ip) Type() string {
	return "env.ip"
}

func (r *Ip) Parse(config config.Exam) (exams.Exam, error) {
	if config.Type != r.Type() {
		return nil, &exams.WrongExamParserError{Source: config.Type, Using: r.Type()}
	}

	if len(config.Vars) == 0 {
		return nil, &VarsUnsetError{Exam: r.Type()}
	}

	return &Ip{config.Vars}, nil
}

// TODO: Refactor this
func (r *Ip) Examinate() (bool, []error) {
	for _, v := range r.Vars {
		ipv4 := &Ipv4{Vars: []string{v}}
		ipv6 := &Ipv6{Vars: []string{v}}

		ok, err := exams.EitherExaminate([]func() (bool, []error){
			ipv4.Examinate,
			ipv6.Examinate,
		})

		if !ok {
			return false, err
		}
	}

	return true, nil
}
