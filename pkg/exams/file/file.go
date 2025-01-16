package file

import (
	"os"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
	"github.com/OJarrisonn/medik/pkg/format"
)

// A report that is returned from a `file.*` exam
type FileReport struct {
	Type     string
	Lvl      int
	Statuses []FileStatus
}

func (r *FileReport) Level() int {
	return r.Lvl
}

func (r *FileReport) Format(verbosity int) (int, string, string) {
	statuses := ""

	for _, status := range r.Statuses {
		if status.Lvl >= verbosity {
			statuses += format.ReportStatus(status.Path, status.Message, status.Lvl) + "\n"
		}
	}

	return r.Lvl, format.ReportHeader(r.Type, r.Lvl), statuses
}

// A status from a part of the execution of a `file.*` exam
type FileStatus struct {
	Lvl     int
	Path    string
	Message string
}

var parsers = map[string]func(config config.Exam) (exams.Exam, error){
	exams.ExamType[*Path]():       exams.ExamParse[*Path](),
	exams.ExamType[*IsFile]():     exams.ExamParse[*IsFile](),
	exams.ExamType[*IsDir]():      exams.ExamParse[*IsDir](),
	exams.ExamType[*IsEmpty]():    exams.ExamParse[*IsEmpty](),
	exams.ExamType[*IsNotEmpty](): exams.ExamParse[*IsNotEmpty](),
}

// Function to get a parser for a given type `env.*`
// The type is the part after `env.` in the exams.Exam type
// Returns the parser and a boolean indicating if the parser was found
func GetParser(ty string) (func(config config.Exam) (exams.Exam, error), bool) {
	if parser, ok := parsers[ty]; ok {
		return parser, ok
	}

	return nil, false
}

// Default implementation for Examinate method of exams.Exam. It checks for the existence of the environment
// variables in `vars`. For those who exist, it validates the value using the `validate` function which should
// return a boolean (valid or not) and an error if not valid. Those who are not set are considered invalid and
// append an UnsetEnvVarError to the errors slice. If no errors are found, it returns true and nil.
func DefaultExaminate(exam string, logLevel int, paths []string, validate func(path string, stat os.FileInfo) FileStatus) *FileReport {
	statuses := []FileStatus{}
	level := 0

	for _, path := range paths {
		stat, err := os.Stat(path)
		if err != nil {
			level = logLevel
			statuses = append(statuses, inexistentPathStatus(path, logLevel))
		} else {
			status := validate(path, stat)

			if status.Lvl > logLevel {
				status.Lvl = logLevel
			}

			if status.Lvl > level {
				level = status.Lvl
			}

			statuses = append(statuses, status)
		}
	}

	return &FileReport{Type: exam, Lvl: level, Statuses: statuses}
}

func DefaultParse[E exams.Exam](config config.Exam, f func(config config.Exam) (exams.Exam, error)) (exams.Exam, error) {
	var e E
	ty := e.Type()
	if config.Type != ty {
		return nil, &exams.WrongExamParserError{Source: config.Type, Using: ty}
	}

	if len(config.Paths) == 0 {
		return nil, &exams.MissingFieldError{Field: "paths", Exam: ty}
	}

	return f(config)
}

func inexistentPathStatus(name string, level int) FileStatus {
	return FileStatus{
		Lvl:     level,
		Path:    name,
		Message: "does not exist",
	}
}

func validPathStatus(name string) FileStatus {
	return FileStatus{
		Lvl:     0,
		Path:    name,
		Message: "is valid",
	}
}

func invalidPathStatus(name string, level int, message string) FileStatus {
	return FileStatus{
		Lvl:     level,
		Path:    name,
		Message: message,
	}
}
