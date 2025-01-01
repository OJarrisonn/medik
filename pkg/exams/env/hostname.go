package env

import (
	"fmt"
	neturl "net/url"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
	"github.com/OJarrisonn/medik/pkg/medik"
)

// A rule to check if an environment variable is set to a hostname
//
// type: env.hostname,
// vars: []string,
// protocol: string
type Hostname struct {
	Vars     []string
	Level    int
	Protocol string
}

func (r *Hostname) Type() string {
	return "env.hostname"
}

func (r *Hostname) Parse(conf config.Exam) (exams.Exam, error) {
	return DefaultParse[*Hostname](conf, func(conf config.Exam) (exams.Exam, error) {
		return &Hostname{conf.Vars, medik.LogLevelFromStr(conf.Level), conf.Protocol}, nil
	})
}

func (r *Hostname) Examinate() exams.Report {
	return DefaultExaminate(r.Type(), r.Level, r.Vars, func(name, value string) EnvStatus {
		ok, _ := r.validateUrl(value)

		if !ok {
			return invalidEnvVarStatus(name, r.Level, value, "value should be a valid URL")
		}

		return validEnvVarStatus(name)
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
