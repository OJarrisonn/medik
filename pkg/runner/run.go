package runner

import (
	"fmt"
	"slices"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/exams"
	"github.com/OJarrisonn/medik/pkg/parse"
)

type UnknownExamError struct {
	ExamType string
}

func (e *UnknownExamError) Error() string {
	return fmt.Sprintf("unknown exam: %v", e.ExamType)
}

func Run(config *config.Medik, protocols []string) (bool, []exams.Report, error) {
	vitalsSuccess, vitalsReports, vitalsError := runVitals(config.Vitals)

	if vitalsError != nil {
		return false, nil, vitalsError
	}

	_, checksReports, checksError := runChecks(config.Checks)

	if checksError != nil {
		return false, nil, checksError
	}

	protocolsUsed := config.Protocols

	for k := range protocolsUsed {
		if !slices.Contains(protocols, k) {
			delete(protocolsUsed, k)
		}
	}

	protocolsSuccess, protocolsReports, protocolsError := runProtocols(protocolsUsed)

	if protocolsError != nil {
		return false, nil, protocolsError
	}

	success := vitalsSuccess && protocolsSuccess

	return success, append(append(vitalsReports, checksReports...), protocolsReports...), nil
}

func runVitals(vitals []config.Exam) (bool, []exams.Report, error) {
	success := true
	reports := []exams.Report{}

	for _, v := range vitals {
		if parse, ok := parse.GetExamParser(v.Type); !ok {
			return false, nil, &UnknownExamError{ExamType: v.Type}
		} else {
			exam, err := parse(v)
			if err != nil {
				return false, nil, err
			}

			report := exam.Examinate()
			if !report.Succeed() {
				success = false
			}

			reports = append(reports, report)
		}
	}

	return success, reports, nil
}

func runChecks(checks []config.Exam) (bool, []exams.Report, error) {
	success := true
	reports := []exams.Report{}

	for _, c := range checks {
		if parse, ok := parse.GetExamParser(c.Type); !ok {
			return false, nil, &UnknownExamError{ExamType: c.Type}
		} else {
			exam, err := parse(c)
			if err != nil {
				return false, nil, err
			}

			report := exam.Examinate()
			if !report.Succeed() {
				success = false
			}

			reports = append(reports, report)
		}
	}

	return success, reports, nil
}

func runProtocols(protocols map[string]config.Protocol) (bool, []exams.Report, error) {
	reports := []exams.Report{}
	success := true

	for _, p := range protocols {
		vitalsSuccess, vitalsReports, vitalsError := runVitals(p.Vitals)

		if vitalsError != nil {
			return false, nil, vitalsError
		}

		reports = append(reports, vitalsReports...)

		if !vitalsSuccess {
			success = false
		}

		_, checksReports, checksError := runChecks(p.Checks)

		if checksError != nil {
			return false, nil, checksError
		}

		reports = append(reports, checksReports...)
	}

	return success, reports, nil
}
