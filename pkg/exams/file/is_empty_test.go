package file

import (
	"testing"

	"github.com/OJarrisonn/medik/pkg/config"
)

func TestIsEmpty_Type(t *testing.T) {
	i := &IsEmpty{}
	if got, want := i.Type(), "file.is-empty"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestIsEmpty_Parse(t *testing.T) {
	i := &IsEmpty{}
	if got, _ := i.Parse(config.Exam{}); got != nil {
		t.Errorf("got %v, want nil", got)
	}
}
