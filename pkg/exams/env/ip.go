package env

import (
	"fmt"

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
		return nil, fmt.Errorf("vars is not set for env.ip")
	}

	return &Ip{config.Vars}, nil
}

// TODO: Refactor this
func (r *Ip) Examinate() (bool, error) {
	unset := []string{}
	invalid := []string{}

	for _, v := range r.Vars {
		ipv4 := &Ipv4{Vars: []string{v}}
		ipv6 := &Ipv6{Vars: []string{v}}

		if ok, _ := ipv4.Examinate(); !ok {
			if ok, _ := ipv6.Examinate(); !ok {
				invalid = append(invalid, v)
			}
		}
	}

	err := ""

	if len(unset) > 0 {
		err += fmt.Sprintf("environment variables not set %v\n", unset)
	}

	if len(invalid) > 0 {
		err += fmt.Sprintf("environment variables not set to valid IP addresses %v\n", invalid)
	}

	if err != "" {
		return false, fmt.Errorf("%v", err)
	}

	return true, nil
}
