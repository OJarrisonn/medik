// This package defines rules behaviour for a Medik configuration file
package exams

// Interface that describes a rule
type Exam interface {
	// Examinate checks if a rule is being enforced
	// Returns true if the rule is being enforced, false otherwise
	// Returns an error if any underlying operation fails or the rule is not being enforced
	Examinate() (bool, error)
}
