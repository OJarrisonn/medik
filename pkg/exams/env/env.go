package env

import (
	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
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
