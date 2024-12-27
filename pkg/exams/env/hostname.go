package env

import (
	"fmt"
	neturl "net/url"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
)

// A rule to check if an environment variable is set to a hostname
//
// type: env.hostname,
// vars: []string,
// protocol: string
type Hostname struct {
	Vars     []string
	Protocol string
}

func (r *Hostname) Type() string {
	return "env.hostname"
}

func (r *Hostname) Parse(config config.Exam) (exams.Exam, error) {
	if config.Type != r.Type() {
		return nil, &exams.WrongExamParserError{Source: config.Type, Using: r.Type()}
	}

	if len(config.Vars) == 0 {
		return nil, &VarsUnsetError{Exam: r.Type()}
	}

	return &Hostname{config.Vars, config.Protocol}, nil
}

func (r *Hostname) Examinate() (bool, []error) {
	return DefaultExaminate(r.Vars, func(name, value string) (bool, error) {
		ok, _ := r.validateUrl(value)

		if !ok {
			return false, &InvalidEnvVarError{Var: name, Value: value, Message: "value should be a valid URL"}
		}

		return true, nil
	})
}

func (r *Hostname) validateUrl(rawUrl string) (bool, error) {
	url, err := neturl.Parse(rawUrl)
	if err != nil {
		return false, err
	}

	if r.Protocol == "" {
		return true, nil
	}

	if url.Scheme != r.Protocol {
		return false, fmt.Errorf("URL %v does not match protocol %v", rawUrl, r.Protocol)
	}

	return true, nil
}
