package env

import (
	"testing"

	"github.com/OJarrisonn/medik/pkg/medik"
	"github.com/stretchr/testify/assert"
)

func TestEnvIsSet(t *testing.T) {
	exam := &IsSet{Vars: []string{"VAR1", "VAR2"}}

	// Test when environment variables are not set
	report := exam.Examinate()
	ok, header, body := report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables are set
	t.Setenv("VAR1", "value1")
	t.Setenv("VAR2", "value2")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.OK, ok)
	assert.NotEmpty(t, header)
	assert.Empty(t, body)
}
