package env

import (
	"testing"

	"github.com/OJarrisonn/medik/pkg/medik"
	"github.com/stretchr/testify/assert"
)

func TestEnvDirExists(t *testing.T) {
	exam := &Dir{Vars: []string{"VAR1", "VAR2"}, Exists: true}

	// Test when environment variables are not set
	report := exam.Examinate()
	ok, header, body := report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables are not valid directories
	t.Setenv("VAR1", "/invalid/path")
	t.Setenv("VAR2", "/etc")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables are valid directories
	t.Setenv("VAR1", "/etc")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.OK, ok)
	assert.NotEmpty(t, header)
	assert.Empty(t, body)
}

func TestEnvDirNotExists(t *testing.T) {
	exam := &Dir{Vars: []string{"VAR1", "VAR2"}, Exists: false}

	// Test when environment variables are not set
	report := exam.Examinate()
	ok, header, body := report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables are valid directories
	t.Setenv("VAR1", "/etc")
	t.Setenv("VAR2", "/invalid/path")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables are not valid directories
	t.Setenv("VAR1", "/invalid/path")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.OK, ok)
	assert.NotEmpty(t, header)
	assert.Empty(t, body)
}
