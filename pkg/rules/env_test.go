package rules

import (
	"testing"
)

func TestEnvIsSet(t *testing.T) {
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
		t.Setenv(k, v)
	}

	for k := range set_vars {
		ok, err := (&EnvIsSet{EnvVar: k}).Validate()

		if !ok {
			t.Errorf("%s is not set\n", k)
		}

		if err != nil {
			t.Errorf("%v\n", err)
		}
	}

	for _, v := range unset_vars {
		ok, err := (&EnvIsSet{EnvVar: v}).Validate()

		if ok {
			t.Errorf("%s is set\n", v)
		}

		if err == nil {
			t.Errorf("No error was raised\n")
		} else {
			t.Logf("%v\n", err)
		}
	}
}

func TestEnvRegexPass(t *testing.T) {
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
		t.Setenv(k, v)
	}

	for k := range set_vars {
		if ok, _ := (&EnvRegex{k, regex}).Validate(); !ok {
			t.Errorf("%s is not set\n", k)
		}
	}

	for _, v := range unset_vars {
		if ok, err := (&EnvRegex{v, regex}).Validate(); ok {
			t.Errorf("%s is set\n", v)
		} else {
			t.Logf("%v\n", err)
		}
	}
}

func TestEnvRegexFail(t *testing.T) {
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

	regex := "bar[0-9]"

	for k, v := range set_vars {
		t.Setenv(k, v)
	}

	for k := range set_vars {
		switch _, err := (&EnvRegex{k, regex}).Validate(); err.(type) {
		case *EnvRegexError:
			t.Log(err)
		default:
			t.Errorf("%s raised an error %v\n", k, err)
		}
	}

	for _, v := range unset_vars {
		if ok, _ := (&EnvRegex{v, regex}).Validate(); ok {
			t.Errorf("%s is set\n", v)
		}
	}
}

func TestEnvRegexCompileError(t *testing.T) {
	t.Setenv("MEDIK_FOO1", "foo1")
	ok, err := (&EnvRegex{"MEDIK_FOO1", "foo[0-9"}).Validate()

	if ok || err == nil {
		t.Errorf("MEDIK_FOO1 is set and the regex `foo[0-9` was approved\n")
	}

	t.Logf("%v %T", err, err)
}

func TestEnvOptionPass(t *testing.T) {
	set_vars := map[string]string{
		"MEDIK_FOO1": "foo1",
		"MEDIK_FOO2": "foo2",
		"MEDIK_FOO3": "foo3",
	}

	options := []string{"foo1", "foo2", "foo3"}

	for k, v := range set_vars {
		t.Setenv(k, v)
	}

	for k := range set_vars {
		if ok, _ := (&EnvOption{k, options}).Validate(); !ok {
			t.Errorf("%s is not set\n", k)
		}
	}
}

func TestEnvOptionUnset(t *testing.T) {
	unset_vars := []string{
		"MEDIK_FOO4",
		"MEDIK_FOO5",
		"MEDIK_FOO6",
	}

	options := []string{"foo1", "foo2", "foo3"}

	for _, v := range unset_vars {
		if ok, _ := (&EnvOption{v, options}).Validate(); ok {
			t.Errorf("%s is set\n", v)
		}
	}
}

func TestEnvOptionFail(t *testing.T) {
	set_vars := map[string]string{
		"MEDIK_FOO1": "foo4",
		"MEDIK_FOO2": "foo5",
		"MEDIK_FOO3": "foo3",
	}

	options := []string{"foo1", "foo2"}

	for k, v := range set_vars {
		t.Setenv(k, v)
	}

	for k := range set_vars {
		if ok, err := (&EnvOption{k, options}).Validate(); ok {
			t.Errorf("%s is valid\n", k)
		} else {
			t.Logf("%v\n", err)
		}
	}
}

func TestEnvIntegerPass(t *testing.T) {
	set_vars := map[string]string{
		"MEDIK_FOO1": "1",
		"MEDIK_FOO2": "-2",
		"MEDIK_FOO3": "3",
	}

	for k, v := range set_vars {
		t.Setenv(k, v)
	}

	for k := range set_vars {
		if ok, err := (&EnvInteger{k}).Validate(); !ok {
			t.Errorf("%s is not set to an integer, %v\n", k, err)
		}
	}
}

func TestEnvIntegerFail(t *testing.T) {
	set_vars := map[string]string{
		"MEDIK_FOO1": ".1",
		"MEDIK_FOO2": "2O",
		"MEDIK_FOO3": "3f",
	}

	for k, v := range set_vars {
		t.Setenv(k, v)
	}

	for k := range set_vars {
		if ok, err := (&EnvInteger{k}).Validate(); ok {
			t.Errorf("%s is being accepted as an integer integer\n", k)
		} else {
			t.Logf("%v\n", err)
		}
	}
}

func TestEnvIntegerUnset(t *testing.T) {
	unset_vars := []string{
		"MEDIK_FOO4",
		"MEDIK_FOO5",
		"MEDIK_FOO6",
	}

	for _, v := range unset_vars {
		if ok, _ := (&EnvInteger{v}).Validate(); ok {
			t.Errorf("%s is set\n", v)
		}
	}
}

func TestEnvFloatPass(t *testing.T) {
	set_vars := map[string]string{
		"MEDIK_FOO1": "1.0",
		"MEDIK_FOO2": "2.0",
		"MEDIK_FOO3": "-3.0",
	}

	for k, v := range set_vars {
		t.Setenv(k, v)
	}

	for k := range set_vars {
		if ok, err := (&EnvFloat{k}).Validate(); !ok {
			t.Errorf("%s is not set to a float, %v\n", k, err)
		}
	}
}

func TestEnvFloatFail(t *testing.T) {
	set_vars := map[string]string{
		"MEDIK_FOO1": "x.1",
		"MEDIK_FOO2": "2O",
		"MEDIK_FOO3": "3f",
	}

	for k, v := range set_vars {
		t.Setenv(k, v)
	}

	for k := range set_vars {
		if ok, err := (&EnvFloat{k}).Validate(); ok {
			t.Errorf("%s is being accepted as a float\n", k)
		} else {
			t.Logf("%v\n", err)
		}
	}
}

func TestEnvFloatUnset(t *testing.T) {
	unset_vars := []string{
		"MEDIK_FOO4",
		"MEDIK_FOO5",
		"MEDIK_FOO6",
	}

	for _, v := range unset_vars {
		if ok, _ := (&EnvFloat{v}).Validate(); ok {
			t.Errorf("%s is set\n", v)
		}
	}
}

