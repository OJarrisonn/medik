package env

import (
	"fmt"
	"os"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
	"github.com/OJarrisonn/medik/pkg/format"
	"github.com/OJarrisonn/medik/pkg/medik"
)

// Function to get a parser for a given type `env.*`
// The type is the part after `env.` in the exams.Exam type
// Returns the parser and a boolean indicating if the parser was found
func GetParser(ty string) (func(config config.Exam) (exams.Exam, error), bool) {
	if parser, ok := parsers[ty]; ok {
		return parser, ok
	}

	return nil, false
}

var parsers = map[string]func(config config.Exam) (exams.Exam, error){
	exams.ExamType[*IsSet]():      exams.ExamParse[*IsSet](),
	exams.ExamType[*NotEmpty]():   exams.ExamParse[*NotEmpty](),
	exams.ExamType[*Regex]():      exams.ExamParse[*Regex](),
	exams.ExamType[*Option]():     exams.ExamParse[*Option](),
	exams.ExamType[*Int]():        exams.ExamParse[*Int](),
	exams.ExamType[*IntRange]():   exams.ExamParse[*IntRange](),
	exams.ExamType[*Float]():      exams.ExamParse[*Float](),
	exams.ExamType[*FloatRange](): exams.ExamParse[*FloatRange](),
	exams.ExamType[*File]():       exams.ExamParse[*File](),
	exams.ExamType[*Dir]():        exams.ExamParse[*Dir](),
	exams.ExamType[*Ipv4]():       exams.ExamParse[*Ipv4](),
	exams.ExamType[*Ipv6]():       exams.ExamParse[*Ipv6](),
	exams.ExamType[*Ip]():         exams.ExamParse[*Ip](),
	exams.ExamType[*Hostname]():   exams.ExamParse[*Hostname](),
}

type EnvReport struct {
	Type     string
	Lvl      int
	Statuses []EnvStatus
}

type EnvStatus struct {
	Lvl     int
	Var     string
	Message string
}

func (r *EnvReport) Level() int {
	return r.Lvl
}

func (r *EnvReport) Format(verbosity int) (int, string, string) {
	statuses := ""

	for _, status := range r.Statuses {
		if status.Lvl >= verbosity {
			statuses += format.ReportStatus(status.Var, status.Message, status.Lvl) + "\n"
		}
	}

	return r.Lvl, format.ReportHeader(r.Type, r.Lvl), statuses
}

type VarsUnsetError struct {
	Exam string
}

func (e *VarsUnsetError) Error() string {
	return "`vars` field is not set for exam " + e.Exam
}

func validEnvVarStatus(name string) EnvStatus {
	return EnvStatus{
		Lvl:     medik.OK,
		Var:     name,
		Message: "is valid",
	}
}

// Function to create a status for an environment variable that is not set
func unsetEnvVarStatus(name string, level int) EnvStatus {
	return EnvStatus{
		Lvl:     level,
		Var:     name,
		Message: "is not set",
	}
}

// Function to create a status for an environment variable whose value is invalid
func invalidEnvVarStatus(name string, level int, value, message string) EnvStatus {
	return EnvStatus{
		Lvl:     level,
		Var:     name,
		Message: fmt.Sprintf("'%v' is not valid: %v", value, message),
	}
}

// Default implementation for Examinate method of exams.Exam. It checks for the existence of the environment
// variables in `vars`. For those who exist, it validates the value using the `validate` function which should
// return a boolean (valid or not) and an error if not valid. Those who are not set are considered invalid and
// append an UnsetEnvVarError to the errors slice. If no errors are found, it returns true and nil.
func DefaultExaminate(exam string, logLevel int, vars []string, validate func(name, value string) EnvStatus) *EnvReport {
	statuses := []EnvStatus{}
	level := 0

	for _, name := range vars {
		value, ok := os.LookupEnv(name)
		if !ok {
			level = logLevel
			statuses = append(statuses, unsetEnvVarStatus(name, logLevel))
		} else {
			status := validate(name, value)

			if status.Lvl > logLevel {
				status.Lvl = logLevel
			}

			if status.Lvl > level {
				level = status.Lvl
			}

			statuses = append(statuses, status)
		}
	}

	return &EnvReport{Type: exam, Lvl: level, Statuses: statuses}
}

func DefaultParse[E exams.Exam](config config.Exam, f func(config config.Exam) (exams.Exam, error)) (exams.Exam, error) {
	var e E
	ty := e.Type()
	if config.Type != ty {
		return nil, &exams.WrongExamParserError{Source: config.Type, Using: ty}
	}

	if len(config.Vars) == 0 {
		return nil, &VarsUnsetError{Exam: ty}
	}

	return f(config)
}
