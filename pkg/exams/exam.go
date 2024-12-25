// This package defines rules behaviour for a Medik configuration file
package exams

import (
	"github.com/OJarrisonn/medik/pkg/config"
)

// Interface that describes a rule
type Exam interface {
	// Type returns the type of the exam
	// This is used to parse the config.Exam by selecting the correct exam parser
	// This method is always called on a zero value of the implementing struct
	Type() string

	// Try parses an []exams.Exam from a config.Exam
	// Returns an error if the config.Exam is invalid
	// This method is always called on a zero value of the implementing struct
	Parse(config config.Exam) (Exam, error)

	// Examinate checks if a rule is being enforced
	// Returns true if the rule is being enforced, false otherwise
	// Returns an error if any underlying operation fails or the rule is not being enforced
	Examinate() (bool, error)
}

// Returns the Parse method from an Exam
// Since this method is decoupled from the values stored in the struct
func ExamParse[E Exam]() func(config config.Exam) (Exam, error) {
	var e E
	return e.Parse
}

// Returns the Type string from an Exam
func ExamType[E Exam]() string {
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
