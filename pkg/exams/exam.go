// This package defines rules behaviour for a Medik configuration file
package exams

import (
	"strings"

	"github.com/OJarrisonn/medik/pkg/config"
)

type ExamParser func(config config.Exam) (Exam, error)

// Interface that describes a rule
type Exam interface {
	// Type returns the type of the exam
	// This is used to parse the config.Exam by selecting the correct exam parser
	Type() string

	// Try parses an []exams.Exam from a config.Exam
	// Returns an error if the config.Exam is invalid
	Parse(config config.Exam) (Exam, error)

	// Examinate checks if a rule is being enforced
	// Returns true if the rule is being enforced, false otherwise
	// Returns an error if any underlying operation fails or the rule is not being enforced
	Examinate() (bool, error)
}

// Parser returns an ExamParser for the given type
// A type is a string in the format of "category.kind"
// Returns a bool indicating if the parser was found
func Parser(t string) (ExamParser, bool) {
	cat, _, _ := strings.Cut(t, ".")

	switch cat {
	default:
		return nil, false
	}
}

func GetParserForType[E Exam]() ExamParser {
	var e E
	return e.Parse
}

func GetTypeForType[E Exam]() string {
	var e E
	return e.Type()
}

type WrongExamParserError struct {
	Got,
	Expected string
}

func (e *WrongExamParserError) Error() string {
	return "wrong exam parser called expected " + e.Expected + " got " + e.Got
}
