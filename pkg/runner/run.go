package runner

import (
	"fmt"
	"slices"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
	"github.com/OJarrisonn/medik/pkg/medik"
	"github.com/OJarrisonn/medik/pkg/parse"
)

type UnknownExamError struct {
	ExamType string
}

func (e *UnknownExamError) Error() string {
	return fmt.Sprintf("unknown exam: %v", e.ExamType)
}

func Run(config *config.Medik, protocols []string) (int, []exams.Report, error) {
	examsSuccesses, examsReports, examsError := runExams(config.Exams)

	if examsError != nil {
		return medik.ERROR, nil, examsError
	}

	protocolsUsed := config.Protocols

	for k := range protocolsUsed {
		if !slices.Contains(protocols, k) {
			delete(protocolsUsed, k)
		}
	}

	protocolsSuccess, protocolsReports, protocolsError := runProtocols(protocolsUsed)

	if protocolsError != nil {
		return medik.ERROR, nil, protocolsError
	}

	success := max(examsSuccesses, protocolsSuccess)

	return success, append(examsReports, protocolsReports...), nil
}

func runExams(exs []config.Exam) (int, []exams.Report, error) {
	reports := []exams.Report{}
	success := medik.OK

	for _, v := range exs {
		if parse, ok := parse.GetExamParser(v.Type); !ok {
			return medik.ERROR, nil, &UnknownExamError{ExamType: v.Type}
		} else {
			exam, err := parse(v)
			if err != nil {
				return medik.ERROR, nil, err
			}

			report := exam.Examinate()
			if report.Level() > medik.OK {
				success = report.Level()
			}

			reports = append(reports, report)
		}
	}

	return success, reports, nil
}

func runProtocols(protocols map[string]config.Protocol) (int, []exams.Report, error) {
	reports := []exams.Report{}
	success := medik.OK

	for _, p := range protocols {
		examsSuccess, examsReports, examsError := runExams(p.Exams)

		if examsError != nil {
			return medik.ERROR, nil, examsError
		}

		reports = append(reports, examsReports...)

		if examsSuccess > success {
			success = examsSuccess
		}
	}

	return success, reports, nil
}
