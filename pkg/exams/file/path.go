package file

import (
	"fmt"
	"os"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
	"github.com/OJarrisonn/medik/pkg/medik"
)

type Path struct {
	Paths  []string
	Level  int
	Exists bool
}

// Type returns the type of the exam
// This is used to parse the config.Exam by selecting the correct exam parser
// This method is always called on a zero value of the implementing struct
func (p *Path) Type() string {
	return "file.path"
}

// Try parses an []exams.Exam from a config.Exam
// Returns an error if the config.Exam is invalid
// This method is always called on a zero value of the implementing struct
func (p *Path) Parse(conf config.Exam) (exams.Exam, error) {
	return DefaultParse[*Path](conf, func(config config.Exam) (exams.Exam, error) {
		return &Path{config.Paths, medik.LogLevelFromStr(config.Level), config.Exists}, nil
	})
}

// Examinate checks if a rule is being enforced
// Returns true if the rule is being enforced, false otherwise
// Returns an error if any underlying operation fails or the rule is not being enforced
func (p *Path) Examinate() exams.Report {
	statuses := []FileStatus{}
	level := 0

	for _, path := range p.Paths {
		_, err := os.Stat(path)
		if (err == nil) == p.Exists {
			statuses = append(statuses, validPathStatus(path))
		} else {
			existence := "not "
			if p.Exists {
				existence = ""
			}
			message := fmt.Sprintf("path should %vexist", existence)
			statuses = append(statuses, invalidPathStatus(path, p.Level, message))
			level = p.Level
		}
	}

	return &FileReport{Type: p.Type(), Lvl: level, Statuses: statuses}
}
