package env

import (
	"fmt"
	"regexp"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
	"github.com/OJarrisonn/medik/pkg/medik"
)

// Check if an environment variable is set and matches a regex
// The regex is a string that will be compiled into a regular expression using Go's regexp package
//
// type: env.regex,
// vars: []string,
// regex: string
type Regex struct {
	Vars  []string
	Level int
	Regex *regexp.Regexp
}

func (r *Regex) Type() string {
	return "env.regex"
}

func (r *Regex) Parse(conf config.Exam) (exams.Exam, error) {
	return DefaultParse[*Regex](conf, func(config config.Exam) (exams.Exam, error) {
		if config.Regex == "" {
			return nil, &exams.MissingFieldError{Field: "regex", Exam: r.Type()}
		}

		regexp, rerr := regexp.Compile(config.Regex)

		if rerr != nil {
			return nil, &exams.FieldValueError{Field: "regex", Exam: r.Type(), Value: config.Regex, Message: rerr.Error()}
		}

		return &Regex{config.Vars, medik.LogLevelFromStr(config.Level), regexp}, nil
	})
}

func (r *Regex) Examinate() exams.Report {
	return DefaultExaminate(r.Type(), r.Level, r.Vars, func(name, value string) EnvStatus {
		if !r.Regex.MatchString(value) {
			return invalidEnvVarStatus(name, r.Level, value, r.ErrorMessage())
		}

		return validEnvVarStatus(name)
	})
}

func (r *Regex) ErrorMessage() string {
	return fmt.Sprintf("value should match regex %v", r.Regex)
}
