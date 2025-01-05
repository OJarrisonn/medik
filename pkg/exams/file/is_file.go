package file

import (
	"os"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
	"github.com/OJarrisonn/medik/pkg/medik"
)

type IsFile struct {
	Paths []string
	Level int
}

// Type returns the type of the exam
// This is used to parse the config.Exam by selecting the correct exam parser
// This method is always called on a zero value of the implementing struct
func (i *IsFile) Type() string {
	return "file.is-file"
}

// Try parses an []exams.Exam from a config.Exam
// Returns an error if the config.Exam is invalid
// This method is always called on a zero value of the implementing struct
func (i *IsFile) Parse(conf config.Exam) (exams.Exam, error) {
	return DefaultParse[*IsFile](conf, func(config config.Exam) (exams.Exam, error) {
		return &IsFile{config.Paths, medik.LogLevelFromStr(config.Level)}, nil
	})
}

// Examinate checks if a rule is being enforced
// Returns true if the rule is being enforced, false otherwise
// Returns an error if any underlying operation fails or the rule is not being enforced
func (i *IsFile) Examinate() exams.Report {
	return DefaultExaminate(i.Type(), i.Level, i.Paths, func(path string, stat os.FileInfo) FileStatus {
		if !stat.Mode().IsRegular() {
			return invalidPathStatus(path, i.Level, "path isn't a regular file")
		}

		return validPathStatus(path)
	})
}
