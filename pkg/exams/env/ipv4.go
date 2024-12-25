package env

import (
	"fmt"
	"os"
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
		return nil, &exams.WrongExamParserError{Got: config.Type, Expected: r.Type()}
	}

	if len(config.Vars) == 0 {
		return nil, fmt.Errorf("vars is not set for env.ipv4")
	}

	return &Ipv4{config.Vars}, nil
}

func (r *Ipv4) Examinate() (bool, error) {
	unset := []string{}
	invalid := []string{}
	regexp := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)

	for _, v := range r.Vars {
		val, ok := os.LookupEnv(v)
		if !ok {
			unset = append(unset, v)
		} else if !regexp.MatchString(val) {
			invalid = append(invalid, v)
		}
	}

	err := ""

	if len(unset) > 0 {
		err += fmt.Sprintf("environment variables not set %v\n", unset)
	}

	if len(invalid) > 0 {
		err += fmt.Sprintf("environment variables not set to valid IPv4 addresses %v\n", invalid)
	}

	if err != "" {
		return false, fmt.Errorf("%v", err)
	}

	return true, nil
}
