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

func TestEnvIntegerRangePass(t *testing.T) {
	set_vars := map[string]string{
		"MEDIK_FOO1": "1",
		"MEDIK_FOO2": "2",
		"MEDIK_FOO3": "3",
	}

	for k, v := range set_vars {
		t.Setenv(k, v)
	}

	for k := range set_vars {
		if ok, err := (&EnvIntegerRange{k, 1, 4}).Validate(); !ok {
			t.Errorf("%s is not set to an integer in the range [1,4), %v\n", k, err)
		}
	}
}

func TestEnvIntegerRangeFail(t *testing.T) {
	set_vars := map[string]string{
		"MEDIK_FOO1": "-1",
		"MEDIK_FOO2": "0",
		"MEDIK_FOO3": "3",
	}

	for k, v := range set_vars {
		t.Setenv(k, v)
	}

	for k := range set_vars {
		if ok, err := (&EnvIntegerRange{k, 1, 3}).Validate(); ok {
			t.Errorf("%s is being accepted as an integer in the range [1,3)\n", k)
		} else {
			t.Logf("%v\n", err)
		}
	}
}

func TestEnvIntegerRangeError(t *testing.T) {
	t.Setenv("MEDIK_FOO1", "a")
	ok, err := (&EnvIntegerRange{"MEDIK_FOO1", 1, 0}).Validate()

	if ok || err == nil {
		t.Errorf("MEDIK_FOO1 is set and the range [1,0) was approved\n")
	}

	t.Logf("%v %T", err, err)
}

func TestEnvIntegerRangeUnset(t *testing.T) {
	unset_vars := []string{
		"MEDIK_FOO4",
		"MEDIK_FOO5",
		"MEDIK_FOO6",
	}

	for _, v := range unset_vars {
		if ok, _ := (&EnvIntegerRange{v, 1, 4}).Validate(); ok {
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

func TestEnvFloatRangeUnset(t *testing.T) {
	unset_vars := []string{
		"MEDIK_FOO4",
		"MEDIK_FOO5",
		"MEDIK_FOO6",
	}

	for _, v := range unset_vars {
		if ok, _ := (&EnvFloatRange{v, 1.0, 4.0}).Validate(); ok {
			t.Errorf("%s is set\n", v)
		}
	}
}

func TestEnvFloatRangePass(t *testing.T) {
	set_vars := map[string]string{
		"MEDIK_FOO1": "1.0",
		"MEDIK_FOO2": "2.0",
		"MEDIK_FOO3": "3.0",
	}

	for k, v := range set_vars {
		t.Setenv(k, v)
	}

	for k := range set_vars {
		if ok, err := (&EnvFloatRange{k, 1.0, 3.0}).Validate(); !ok {
			t.Errorf("%s is not set to a float in the range [1.0,3.0], %v\n", k, err)
		}
	}
}

func TestEnvFloatRangeFail(t *testing.T) {
	set_vars := map[string]string{
		"MEDIK_FOO1": ".9",
		"MEDIK_FOO2": "2.1",
	}

	for k, v := range set_vars {
		t.Setenv(k, v)
	}

	for k := range set_vars {
		if ok, err := (&EnvFloatRange{k, 1.0, 2.0}).Validate(); ok {
			t.Errorf("%s is being accepted as a float in the range [1.0,2.0]\n", k)
		} else {
			t.Logf("%v\n", err)
		}
	}
}

func TestEnvFloatRangeError(t *testing.T) {
	t.Setenv("MEDIK_FOO1", "a.1")
	ok, err := (&EnvFloatRange{"MEDIK_FOO1", 1.0, 0.0}).Validate()

	if ok || err == nil {
		t.Errorf("MEDIK_FOO1 is set and the range [1.0,0.0) was approved\n")
	}

	t.Logf("%v %T", err, err)
}

func TestEnvFloatRangeErrorUnset(t *testing.T) {
	unset_vars := []string{
		"MEDIK_FOO4",
		"MEDIK_FOO5",
		"MEDIK_FOO6",
	}

	for _, v := range unset_vars {
		if ok, _ := (&EnvFloatRange{v, 1.0, 4.0}).Validate(); ok {
			t.Errorf("%s is set\n", v)
		}
	}
}

func TestEnvFilePass(t *testing.T) {
	t.Setenv("MEDIK_FOO1", "/etc/passwd")

	if ok, err := (&EnvFile{"MEDIK_FOO1"}).Validate(); !ok {
		t.Errorf("MEDIK_FOO1 is not set to a file, %v\n", err)
	}
}

func TestEnvFileFail(t *testing.T) {
	t.Setenv("MEDIK_FOO1", "/this/file/does/not/exist/at/all/dingle/bell/beep/boop/foo/bar/baz")

	if ok, err := (&EnvFile{"MEDIK_FOO1"}).Validate(); ok {
		t.Errorf("MEDIK_FOO1 is being accepted as a file\n")
	} else {
		t.Logf("%v\n", err)
	}
}

func TestEnvFileUnset(t *testing.T) {
	if ok, _ := (&EnvFile{"MEDIK_FOO4"}).Validate(); ok {
		t.Errorf("MEDIK_FOO4 is set\n")
	}
}

func TestEnvDirPass(t *testing.T) {
	t.Setenv("MEDIK_FOO1", "/etc")

	if ok, err := (&EnvDir{"MEDIK_FOO1"}).Validate(); !ok {
		t.Errorf("MEDIK_FOO1 is not set to a directory, %v\n", err)
	}
}

func TestEnvDirFail(t *testing.T) {
	t.Setenv("MEDIK_FOO1", "/this/directory/does/not/exist/at/all/dingle/bell/beep/boop/foo/bar/baz")

	if ok, err := (&EnvDir{"MEDIK_FOO1"}).Validate(); ok {
		t.Errorf("MEDIK_FOO1 is being accepted as a directory\n")
	} else {
		t.Logf("%v\n", err)
	}
}

func TestEnvDirNotADir(t *testing.T) {
	t.Setenv("MEDIK_FOO1", "/etc/passwd")

	if ok, err := (&EnvDir{"MEDIK_FOO1"}).Validate(); ok {
		t.Errorf("MEDIK_FOO1 is being accepted as a directory\n")
	} else {
		t.Logf("%v\n", err)
	}
}

func TestEnvDirUnset(t *testing.T) {
	if ok, _ := (&EnvDir{"MEDIK_FOO4"}).Validate(); ok {
		t.Errorf("MEDIK_FOO4 is set\n")
	}
}