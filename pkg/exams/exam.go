// This package defines rules behaviour for a Medik configuration file
package exams

import "github.com/OJarrisonn/medik/pkg/config"

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
