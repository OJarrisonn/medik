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
	Examinate() Report
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

type Report interface {
	// Returns the level of the report
	Level() int

	// Format the report to a printable string
	// Verbose indicates if non-error messages should be included
	// Returns the report level, a string with the report header and a string with the report body
	Format(verbosity int) (int, string, string)
}

// An error to describe a strange scenario where the wrong exam parser was called
// The parser to be called is identified by Type field of an Exam, this becomes the Source field
// The Using field is the Type of the Exam that was called
type WrongExamParserError struct {
	Source,
	Using string
}

func (e *WrongExamParserError) Error() string {
	return "wrong exam parser: using " + e.Using + " parser for a " + e.Source + " exam"
}

// An error to describe a missing field in a config.Exam during its parsing
type MissingFieldError struct {
	Field,
	Exam string
}

func (e *MissingFieldError) Error() string {
	return "missing field `" + e.Field + "` in exam " + e.Exam
}

type FieldValueError struct {
	Field,
	Exam,
	Value,
	Message string
}

func (e *FieldValueError) Error() string {
	return "invalid value '" + e.Value + "' for field `" + e.Field + "` in exam " + e.Exam + ": " + e.Message
}
