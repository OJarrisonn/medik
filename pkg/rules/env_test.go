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
		if ok, err := (&EnvRegex{k, regex}).Validate(); ok {
			t.Errorf("%s was accepted %v\n", k, err)
		} else {
			t.Log(err)
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

func TestEnvIpPass(t *testing.T) {
	// IPv4 and IPv6 addresses
	set_vars := map[string]string{
		"MEDIK_FOO1": "0.0.0.0",
		"MEDIK_FOO2": "::1",
	}

	for k, v := range set_vars {
		t.Setenv(k, v)
	}

	for k := range set_vars {
		if ok, err := (&EnvIpAddr{k}).Validate(); !ok {
			t.Errorf("%s is not set to an IP address, %v\n", k, err)
		}
	}

	if _, err := (&EnvIpv4Addr{"MEDIK_FOO1"}).Validate(); err != nil {
		t.Errorf("MEDIK_FOO1 is not an IPv4 address %v\n", err)
	}

	if _, err := (&EnvIpv6Addr{"MEDIK_FOO2"}).Validate(); err != nil {
		t.Errorf("MEDIK_FOO2 is not an IPv6 address %v\n", err)
	}
}

func TestEnvIpFail(t *testing.T) {
	set_vars := map[string]string{
		"MEDIK_FOO1": "127.o.o.1",
		"MEDIK_FOO2": "abcd:efgh:ijkl:mnop:qrst:uvwx:yzab:cdef",
	}

	for k, v := range set_vars {
		t.Setenv(k, v)
	}

	for k := range set_vars {
		if ok, err := (&EnvIpAddr{k}).Validate(); ok {
			t.Errorf("%s is being accepted as an IP address\n", k)
		} else {
			t.Logf("%v\n", err)
		}
	}

	if _, err := (&EnvIpv4Addr{"MEDIK_FOO1"}).Validate(); err == nil {
		t.Errorf("MEDIK_FOO1 is an IPv4 address\n")
	} else {
		t.Logf("%v\n", err)
	}

	if _, err := (&EnvIpv6Addr{"MEDIK_FOO2"}).Validate(); err == nil {
		t.Errorf("MEDIK_FOO2 is an IPv6 address\n")
	} else {
		t.Logf("%v\n", err)
	}
}

func TestEnvIpUnset(t *testing.T) {
	unset_vars := []string{
		"MEDIK_FOO4",
		"MEDIK_FOO5",
		"MEDIK_FOO6",
	}

	for _, v := range unset_vars {
		if ok, _ := (&EnvIpAddr{v}).Validate(); ok {
			t.Errorf("%s is set\n", v)
		}
	}
}

func TestEnvHostnamePass(t *testing.T) {
	t.Setenv("MEDIK_FOO1", "localhost")
	t.Setenv("MEDIK_FOO2", "http://localhost")

	if ok, err := (&EnvHostname{"MEDIK_FOO1", false, ""}).Validate(); !ok {
		t.Errorf("MEDIK_FOO1 is not set to a valid host, %v\n", err)
	}

	if ok, err := (&EnvHostname{"MEDIK_FOO2", true, "http"}).Validate(); !ok {
		t.Errorf("MEDIK_FOO2 is not set to a valid host, %v\n", err)
	}
}

func TestEnvHostnameFail(t *testing.T) {
	t.Setenv("MEDIK_FOO1", "http://localhost")
	t.Setenv("MEDIK_FOO2", "\n")

	if ok, err := (&EnvHostname{"MEDIK_FOO1", true, "https"}).Validate(); ok {
		t.Errorf("MEDIK_FOO1 is being accepted as a hostname\n")
	} else {
		t.Logf("%v\n", err)
	}

	if ok, err := (&EnvHostname{"MEDIK_FOO2", false, ""}).Validate(); ok {
		t.Errorf("MEDIK_FOO2 is being accepted as a hostname\n")
	} else {
		t.Logf("%v\n", err)
	}
}

func TestEnvHostnameUnset(t *testing.T) {
	unset_vars := []string{
		"MEDIK_FOO4",
		"MEDIK_FOO5",
		"MEDIK_FOO6",
	}

	for _, v := range unset_vars {
		if ok, _ := (&EnvHostname{v, false, ""}).Validate(); ok {
			t.Errorf("%s is set\n", v)
		}
	}
}