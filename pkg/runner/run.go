package runner

import (
	"fmt"
	"slices"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/parse"
)

func Run(config *config.Medik, protocols []string) (bool, []error) {
	vitalsSuccess, vitalsErrors := runVitals(config.Vitals)
	checksSuccess, checksErrors := runChecks(config.Checks)
	protocolsUsed := config.Protocols

	for k := range protocolsUsed {
		if !slices.Contains(protocols, k) {
			delete(protocolsUsed, k)
		}
	}

	protocolsSuccess, protocolsErrors := runProtocols(protocolsUsed)

	success := vitalsSuccess && checksSuccess && protocolsSuccess
	errors := append(vitalsErrors, checksErrors...)
	errors = append(errors, protocolsErrors...)

	if len(errors) == 0 {
		return success, nil
	}

	return success, errors
}

func runVitals(vitals []config.Exam) (bool, []error) {
	errors := []error{}
	success := true

	for _, v := range vitals {
		if parse, ok := parse.GetExamParser(v.Type); !ok {
			errors = append(errors, fmt.Errorf("no parser found for type: %v", v.Type))
			success = false
			continue
		} else {
			exam, err := parse(v)
			if err != nil {
				errors = append(errors, err)
				success = false
				continue
			}

			if ok, errs := exam.Examinate(); !ok {
				errors = append(errors, errs...)
				success = false
			}
		}
	}

	if len(errors) == 0 {
		return true, nil
	}

	return success, errors
}

func runChecks(checks []config.Exam) (bool, []error) {
	errors := []error{}
	success := true

	for _, c := range checks {
		if parse, ok := parse.GetExamParser(c.Type); !ok {
			errors = append(errors, fmt.Errorf("no parser found for type: %v", c.Type))
			success = false
			continue
		} else {
			exam, err := parse(c)
			if err != nil {
				errors = append(errors, err)
				continue
			}

			if ok, errs := exam.Examinate(); !ok {
				errors = append(errors, errs...)
			}
		}
	}

	if len(errors) == 0 {
		return true, nil
	}

	return success, errors
}

func runProtocols(protocols map[string]config.Protocol) (bool, []error) {
	errors := []error{}
	success := true

	for _, p := range protocols {
		vitalsSuccess, vitalsErrors := runVitals(p.Vitals)
		checksSuccess, checksErrors := runChecks(p.Checks)

		if !vitalsSuccess {
			errors = append(errors, vitalsErrors...)
			success = false
		}

		if !checksSuccess {
			errors = append(errors, checksErrors...)
			success = false
		}
	}

	if len(errors) == 0 {
		return true, nil
	}

	return success, errors
}
