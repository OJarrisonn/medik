package env

import (
	"github.com/OJarrisonn/medik/pkg/exams"
)

// Function to get a parser for a given type `env.*`
// The type is the part after `env.` in the exams.Exam type
// Returns the parser and a boolean indicating if the parser was found
func EnvParser(ty string) (exams.ExamParser, bool) {
	if parser, ok := parsers[ty]; ok {
		return parser, ok
	}

	return nil, false
}

var parsers = map[string]exams.ExamParser{
	exams.GetTypeForType[*IsSet]():      exams.GetParserForType[*IsSet](),
	exams.GetTypeForType[*NotEmpty]():   exams.GetParserForType[*NotEmpty](),
	exams.GetTypeForType[*Regex]():      exams.GetParserForType[*Regex](),
	exams.GetTypeForType[*Option]():     exams.GetParserForType[*Option](),
	exams.GetTypeForType[*Int]():        exams.GetParserForType[*Int](),
	exams.GetTypeForType[*IntRange]():   exams.GetParserForType[*IntRange](),
	exams.GetTypeForType[*Float]():      exams.GetParserForType[*Float](),
	exams.GetTypeForType[*FloatRange](): exams.GetParserForType[*FloatRange](),
	exams.GetTypeForType[*File]():       exams.GetParserForType[*File](),
	exams.GetTypeForType[*Dir]():        exams.GetParserForType[*Dir](),
	exams.GetTypeForType[*Ipv4]():       exams.GetParserForType[*Ipv4](),
	exams.GetTypeForType[*Ipv6]():       exams.GetParserForType[*Ipv6](),
	exams.GetTypeForType[*Ip]():         exams.GetParserForType[*Ip](),
	exams.GetTypeForType[*Hostname]():   exams.GetParserForType[*Hostname](),
}
