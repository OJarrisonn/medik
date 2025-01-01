package format

import (
	"fmt"

	"github.com/OJarrisonn/medik/pkg/medik"
)

func ReportHeader(header string, level int) string {
	if medik.NoColor {
		return reportHeaderNoColor(header, level)
	}
	return reportHeaderColor(header, level)
}

func ReportStatus(key, message string, level int) string {
	if medik.NoColor {
		return fmt.Sprintf("\t %s  %s", key, message)
	}

	switch level {
	case medik.OK:
		return medik.SuccessWithBgColor.Sprintf("\t %s ", key) + medik.SuccessColor.Sprintf(" %s", message)
	case medik.WARNING:
		return medik.WarningWithBgColor.Sprintf("\t %s ", key) + medik.WarningColor.Sprintf(" %s", message)
	case medik.ERROR:
		return medik.ErrorWithBgColor.Sprintf("\t %s ", key) + medik.ErrorColor.Sprintf(" %s", message)
	default:
		return fmt.Sprintf("\t %s  %s", key, message)
	}
}

func reportHeaderColor(header string, level int) string {
	title := medik.LogLevel(level)

	switch level {
	case medik.OK:
		return medik.SuccessWithBgColor.Sprintf(" %s ", title) + medik.SuccessColor.Sprintf(" %s", header)
	case medik.WARNING:
		return medik.WarningWithBgColor.Sprintf(" %s ", title) + medik.WarningColor.Sprintf(" %s", header)
	case medik.ERROR:
		return medik.ErrorWithBgColor.Sprintf(" %s ", title) + medik.ErrorColor.Sprintf(" %s", header)
	default:
		return reportHeaderNoColor(header, level)
	}
}

func reportHeaderNoColor(header string, level int) string {
	title := medik.LogLevel(level)
	return fmt.Sprintf(" %s  %s", title, header)
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
