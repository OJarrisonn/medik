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
	Examinate() (bool, []error)
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

// A function that combines multiple validation functions.
// The functions are evaluated in order and short-circuited. It means that if one of the functions returns false,
// the `CompoundExaminate` will return immediately with the errors found until that point.
func CompoundExaminate(validate []func() (bool, []error)) (bool, []error) {
	errors := []error{}
	for _, v := range validate {
		valid, err := v()
		if err != nil {
			errors = append(errors, err...)
		}
		if !valid {
			return false, errors
		}
	}

	if len(errors) == 0 {
		return true, nil
	}

	return true, errors
}

// A function that combines multiple validation functions.
// The functions are evaluated in order and short-circuited. It means that if one of the functions returns true,
// the `CompoundExaminate` will return immediately with the errors found only in the matching validation function.
// If no function returns true, it will return false and all the errors found.
func EitherExaminate(validate []func() (bool, []error)) (bool, []error) {
	errors := []error{}
	for _, v := range validate {
		valid, err := v()
		if valid {
			return true, err
		}
		if err != nil {
			errors = append(errors, err...)
		}
	}

	if len(errors) == 0 {
		return false, nil
	}

	return false, errors
}
