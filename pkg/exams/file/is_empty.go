package file

import (
	"os"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
	"github.com/OJarrisonn/medik/pkg/medik"
)

type IsEmpty struct {
	Paths []string
	Level int
}

// Type returns the type of the exam
// This is used to parse the config.Exam by selecting the correct exam parser
// This method is always called on a zero value of the implementing struct
func (i *IsEmpty) Type() string {
	return "file.is-empty"
}

// Try parses an []exams.Exam from a config.Exam
// Returns an error if the config.Exam is invalid
// This method is always called on a zero value of the implementing struct
func (i *IsEmpty) Parse(conf config.Exam) (exams.Exam, error) {
	return DefaultParse[*IsEmpty](conf, func(config config.Exam) (exams.Exam, error) {
		return &IsEmpty{config.Paths, medik.LogLevelFromStr(config.Level)}, nil
	})
}

// Examinate checks if a rule is being enforced
// Returns true if the rule is being enforced, false otherwise
// Returns an error if any underlying operation fails or the rule is not being enforced
func (i *IsEmpty) Examinate() exams.Report {
	return DefaultExaminate(i.Type(), i.Level, i.Paths, func(path string, stat os.FileInfo) FileStatus {
		if stat.IsDir() {
			empty, err := isDirEmpty(path)
			if err != nil {
				return invalidPathStatus(path, i.Level, err.Error())
			}

			if !empty {
				return invalidPathStatus(path, i.Level, "directory is not empty")
			}

			return validPathStatus(path)
		}

		if stat.Size() != 0 {
			return invalidPathStatus(path, i.Level, "file is not empty")
		}

		return validPathStatus(path)
	})
}

func isDirEmpty(path string) (bool, error) {
	f, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}

	return len(f) == 0, nil
}
