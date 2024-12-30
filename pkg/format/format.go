package format

import (
	"fmt"

	"github.com/OJarrisonn/medik/pkg/medik"
)

func ReportHeader(header string, success string) string {
	if medik.NoColor {
		return reportHeaderNoColor(header, success)
	}
	return reportHeaderColor(header, success)
}

func ReportStatus(key, message string, success bool) string {
	if medik.NoColor {
		if success {
			return fmt.Sprintf("\t %s  %s", key, message)
		} else {
			return fmt.Sprintf("\t %s  %s", key, message)
		}
	}

	if success {
		return medik.SuccessWithBgColor.Sprintf("\t %s ", key) + medik.SuccessColor.Sprintf(" %s", message)
	} else {
		return medik.ErrorWithBgColor.Sprintf("\t %s ", key) + medik.ErrorColor.Sprintf(" %s", message)
	}
}

func reportHeaderColor(header string, success string) string {
	switch success {
	case medik.SUCCESS:
		return medik.SuccessWithBgColor.Sprintf(" %s ", success) + medik.SuccessColor.Sprintf(" %s", header)
	case medik.FAILURE:
		return medik.ErrorWithBgColor.Sprintf(" %s ", success) + medik.ErrorColor.Sprintf(" %s", header)
	default:
		return reportHeaderNoColor(header, success)
	}
}

func reportHeaderNoColor(header string, success string) string {
	return fmt.Sprintf(" %s  %s", success, header)
}

func EnvironmentHealth(healthy bool) string {
	if medik.NoColor {
		if healthy {
			return " Environment Healthy "
		} else {
			return " Environment Unhealthy "
		}
	}

	if healthy {
		return medik.SuccessWithBgColor.Sprintf(" Environment Healthy ")
	}
	return medik.ErrorWithBgColor.Sprintf(" Environment Unhealthy ")
}
