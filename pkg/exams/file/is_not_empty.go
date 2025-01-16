package file

import (
	"os"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
	"github.com/OJarrisonn/medik/pkg/medik"
)

type IsNotEmpty struct {
	Paths []string
	Level int
}

// Type returns the type of the exam
// This is used to parse the config.Exam by selecting the correct exam parser
// This method is always called on a zero value of the implementing struct
func (i *IsNotEmpty) Type() string {
	return "file.is-empty"
}

// Try parses an []exams.Exam from a config.Exam
// Returns an error if the config.Exam is invalid
// This method is always called on a zero value of the implementing struct
func (i *IsNotEmpty) Parse(conf config.Exam) (exams.Exam, error) {
	return DefaultParse[*IsNotEmpty](conf, func(config config.Exam) (exams.Exam, error) {
		return &IsNotEmpty{config.Paths, medik.LogLevelFromStr(config.Level)}, nil
	})
}

// Examinate checks if a rule is being enforced
// Returns true if the rule is being enforced, false otherwise
// Returns an error if any underlying operation fails or the rule is not being enforced
func (i *IsNotEmpty) Examinate() exams.Report {
	return DefaultExaminate(i.Type(), i.Level, i.Paths, func(path string, stat os.FileInfo) FileStatus {
		if stat.IsDir() {
			empty, err := isDirEmpty(path)
			if err != nil {
				return invalidPathStatus(path, i.Level, err.Error())
			}

			if empty {
				return invalidPathStatus(path, i.Level, "directory is empty")
			}

			return validPathStatus(path)
		}

		if stat.Size() == 0 {
			return invalidPathStatus(path, i.Level, "file is empty")
		}

		return validPathStatus(path)
	})
}
