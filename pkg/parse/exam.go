package parse

import (
	"strings"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
	"github.com/OJarrisonn/medik/pkg/exams/env"
	"github.com/OJarrisonn/medik/pkg/exams/file"
)

// Returns the parser for a given type
// A type is a string in the format `category.kind` which identifies which exam will be parsed
// Returns the parser and a boolean indicating if the parser was found
func GetExamParser(ty string) (func(config config.Exam) (exams.Exam, error), bool) {
	category, _, _ := strings.Cut(ty, ".")
	switch category {
	case "env":
		return env.GetParser(ty)
	case "file":
		return file.GetParser(ty)
	default:
		return nil, false
	}
}
