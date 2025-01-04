package format

import (
	"strings"
	"testing"

	"github.com/OJarrisonn/medik/pkg/medik"
)

func TestReportHeaderNoColor(t *testing.T) {
	medik.NoColor = true
	header := "Header"
	for level := range medik.MAX_LEVEL + 2 {
		expected := reportHeaderNoColor(header, level)
		actual := ReportHeader(header, level)
		if actual != expected {
			t.Errorf("Expected %s but got %s", expected, actual)
		}

		if !strings.Contains(actual, header) || !strings.Contains(actual, medik.LogLevel(level)) {
			t.Errorf("Expected %s and %s but got %s", header, medik.LogLevel(level), actual)
		}
	}
}

func TestReportHeaderColor(t *testing.T) {
	medik.NoColor = false
	header := "Header"
	for level := range medik.MAX_LEVEL + 2 {
		expected := reportHeaderColor(header, level)
		actual := ReportHeader(header, level)
		if actual != expected {
			t.Errorf("Expected %s but got %s", expected, actual)
		}

		if !strings.Contains(actual, header) || !strings.Contains(actual, medik.LogLevel(level)) {
			t.Errorf("Expected %s and %s but got %s", header, medik.LogLevel(level), actual)
		}
	}
}

func TestEnvironmentHealthNoColor(t *testing.T) {
	medik.NoColor = true
	for status := range medik.MAX_LEVEL + 2 {
		expected := " Environment Healthy "
		if status >= medik.ERROR {
			expected = " Environment Unhealthy "
		}
		actual := EnvironmentHealth(status)
		if actual != expected {
			t.Errorf("Expected %s but got %s", expected, actual)
		}
	}
}

func TestEnvironmentHealthColor(t *testing.T) {
	medik.NoColor = false
	for status := range medik.MAX_LEVEL + 2 {
		expected := " Environment Healthy "
		if status >= medik.ERROR {
			expected = " Environment Unhealthy "
		}
		actual := EnvironmentHealth(status)
		if actual != expected {
			t.Errorf("Expected %s but got %s", expected, actual)
		}
	}
}

func TestReportStatusNoColor(t *testing.T) {
	medik.NoColor = true
	for level := range medik.MAX_LEVEL + 2 {
		actual := ReportStatus("FOO", "bar baz baz", level)
		if !strings.Contains(actual, "FOO") || !strings.Contains(actual, "bar baz baz") {
			t.Errorf("Expected 'FOO' and 'bar baz baz' but got %s", actual)
		}
	}
}

func TestReportStatusColor(t *testing.T) {
	medik.NoColor = false
	for level := range medik.MAX_LEVEL + 2 {
		actual := ReportStatus("FOO", "bar baz baz", level)
		if !strings.Contains(actual, "FOO") || !strings.Contains(actual, "bar baz baz") {
			t.Errorf("Expected 'FOO' and 'bar baz baz' but got %s", actual)
		}
	}
}
