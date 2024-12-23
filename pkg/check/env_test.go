package check

import (
	"os"
	"testing"
)

func TestCheckEnv(t *testing.T) {
	set_vars := map[string]string{
		"MEDIK_FOO1": "foo1",
		"MEDIK_FOO2": "foo2",
		"MEDIK_FOO3": "foo3",
	}

	unset_vars := []string{
		"MEDIK_FOO4",
		"MEDIK_FOO5",
		"MEDIK_FOO6",
	}

	for k, v := range set_vars {
		os.Setenv(k, v)
	}

	for k := range set_vars {
		if !CheckEnv(k) {
			t.Errorf("%s is not set\n", k)
		}
	}

	for _, v := range unset_vars {
		if CheckEnv(v) {
			t.Errorf("%s is set\n", v)
		}
	}
}

func TestValidateEnvRegex(t *testing.T) {
	set_vars := map[string]string{
		"MEDIK_FOO1": "foo1",
		"MEDIK_FOO2": "foo2",
		"MEDIK_FOO3": "foo3",
	}

	unset_vars := []string{
		"MEDIK_FOO4",
		"MEDIK_FOO5",
		"MEDIK_FOO6",
	}

	regex := "foo[0-9]"

	for k, v := range set_vars {
		os.Setenv(k, v)
	}

	for k := range set_vars {
		if ok, _ := ValidateEnvRegex(k, regex); !ok {
			t.Errorf("%s is not set\n", k)
		}
	}

	for _, v := range unset_vars {
		if ok, _ := ValidateEnvRegex(v, regex); ok {
			t.Errorf("%s is set\n", v)
		}
	}
}