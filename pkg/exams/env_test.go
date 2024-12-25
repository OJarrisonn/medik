package exams

import (
	"regexp"
	"strings"
	"testing"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestEnvParsersCollection(t *testing.T) {
	// Inexistent parser
	parser, ok := envParser("inexistent")

	assert.Nil(t, parser)
	assert.False(t, ok)

	// Existing parsers
	parser, ok = envParser("is-set")

	assert.NotNil(t, parser)
	assert.True(t, ok)
}

func TestEnvParsersContainsAll(t *testing.T) {
	registered := make([]string, len(parsers))

	i := 0
	for k := range parsers {
		registered[i] = k
		i++
	}

	known := []string{
		(&EnvIsSet{}).Type(),
		(&EnvIsSetNotEmpty{}).Type(),
		(&EnvRegex{}).Type(),
		(&EnvOption{}).Type(),
		(&EnvInteger{}).Type(),
		(&EnvIntegerRange{}).Type(),
		(&EnvFloat{}).Type(),
		(&EnvFloatRange{}).Type(),
		(&EnvFile{}).Type(),
		(&EnvDir{}).Type(),
		(&EnvIpv4Addr{}).Type(),
		(&EnvIpv6Addr{}).Type(),
		(&EnvIpAddr{}).Type(),
		(&EnvHostname{}).Type(),
	}

	for i, k := range known {
		known[i], _ = strings.CutPrefix(k, "env.")
	}

	assert.ElementsMatch(t, known, registered)
}

func TestEnvIsSet(t *testing.T) {
	exam := &EnvIsSet{Vars: []string{"VAR1", "VAR2"}}

	// Test when environment variables are not set
	result, err := exam.Examinate()
	assert.False(t, result)
	assert.Error(t, err)

	// Test when environment variables are set
	t.Setenv("VAR1", "value1")
	t.Setenv("VAR2", "value2")
	result, err = exam.Examinate()
	assert.True(t, result)
	assert.NoError(t, err)
}

func TestEnvIsSetNotEmpty(t *testing.T) {
	exam := &EnvIsSetNotEmpty{Vars: []string{"VAR1", "VAR2"}}

	// Test when environment variables are not set
	result, err := exam.Examinate()
	assert.False(t, result)
	assert.Error(t, err)

	// Test when environment variables are set to empty
	t.Setenv("VAR1", "")
	t.Setenv("VAR2", " ")
	result, err = exam.Examinate()
	assert.False(t, result)
	assert.Error(t, err)

	// Test when environment variables are set to non-empty values
	t.Setenv("VAR1", "value1")
	t.Setenv("VAR2", "value2")
	result, err = exam.Examinate()
	assert.True(t, result)
	assert.NoError(t, err)
}

func TestEnvRegex(t *testing.T) {
	regex, _ := regexp.Compile(`^value\d$`)
	exam := &EnvRegex{Vars: []string{"VAR1", "VAR2"}, Regex: regex}

	// Test when environment variables are not set
	result, err := exam.Examinate()
	assert.False(t, result)
	assert.Error(t, err)

	// Test when environment variables do not match regex
	t.Setenv("VAR1", "invalid")
	t.Setenv("VAR2", "value2")
	result, err = exam.Examinate()
	assert.False(t, result)
	assert.Error(t, err)

	// Test when environment variables match regex
	t.Setenv("VAR1", "value1")
	result, err = exam.Examinate()
	assert.True(t, result)
	assert.NoError(t, err)
}

func TestEnvOption(t *testing.T) {
	options := map[string]bool{"option1": true, "option2": true}
	exam := &EnvOption{Vars: []string{"VAR1", "VAR2"}, Options: options}

	// Test when environment variables are not set
	result, err := exam.Examinate()
	assert.False(t, result)
	assert.Error(t, err)

	// Test when environment variables do not match options
	t.Setenv("VAR1", "invalid")
	t.Setenv("VAR2", "option2")
	result, err = exam.Examinate()
	assert.False(t, result)
	assert.Error(t, err)

	// Test when environment variables match options
	t.Setenv("VAR1", "option1")
	result, err = exam.Examinate()
	assert.True(t, result)
	assert.NoError(t, err)
}

func TestEnvInteger(t *testing.T) {
	exam := &EnvInteger{Vars: []string{"VAR1", "VAR2"}}

	// Test when environment variables are not set
	result, err := exam.Examinate()
	assert.False(t, result)
	assert.Error(t, err)

	// Test when environment variables are not integers
	t.Setenv("VAR1", "invalid")
	t.Setenv("VAR2", "123")
	result, err = exam.Examinate()
	assert.False(t, result)
	assert.Error(t, err)

	// Test when environment variables are integers
	t.Setenv("VAR1", "456")
	result, err = exam.Examinate()
	assert.True(t, result)
	assert.NoError(t, err)
}

func TestEnvIntegerRange(t *testing.T) {
	exam := &EnvIntegerRange{Vars: []string{"VAR1", "VAR2"}, Min: 10, Max: 100}

	// Test when environment variables are not set
	result, err := exam.Examinate()
	assert.False(t, result)
	assert.Error(t, err)

	// Test when environment variables are not integers
	t.Setenv("VAR1", "invalid")
	t.Setenv("VAR2", "50")
	result, err = exam.Examinate()
	assert.False(t, result)
	assert.Error(t, err)

	// Test when environment variables are out of range
	t.Setenv("VAR1", "5")
	result, err = exam.Examinate()
	assert.False(t, result)
	assert.Error(t, err)

	// Test when environment variables are within range
	t.Setenv("VAR1", "20")
	result, err = exam.Examinate()
	assert.True(t, result)
	assert.NoError(t, err)
}

func TestEnvFloat(t *testing.T) {
	exam := &EnvFloat{Vars: []string{"VAR1", "VAR2"}}

	// Test when environment variables are not set
	result, err := exam.Examinate()
	assert.False(t, result)
	assert.Error(t, err)

	// Test when environment variables are not floats
	t.Setenv("VAR1", "invalid")
	t.Setenv("VAR2", "123.45")
	result, err = exam.Examinate()
	assert.False(t, result)
	assert.Error(t, err)

	// Test when environment variables are floats
	t.Setenv("VAR1", "456.78")
	result, err = exam.Examinate()
	assert.True(t, result)
	assert.NoError(t, err)
}

func TestEnvFloatRange(t *testing.T) {
	exam := &EnvFloatRange{Vars: []string{"VAR1", "VAR2"}, Min: 10.5, Max: 100.5}

	// Test when environment variables are not set
	result, err := exam.Examinate()
	assert.False(t, result)
	assert.Error(t, err)

	// Test when environment variables are not floats
	t.Setenv("VAR1", "invalid")
	t.Setenv("VAR2", "50.5")
	result, err = exam.Examinate()
	assert.False(t, result)
	assert.Error(t, err)

	// Test when environment variables are out of range
	t.Setenv("VAR1", "5.5")
	result, err = exam.Examinate()
	assert.False(t, result)
	assert.Error(t, err)

	// Test when environment variables are within range
	t.Setenv("VAR1", "20.5")
	result, err = exam.Examinate()
	assert.True(t, result)
	assert.NoError(t, err)
}

func TestEnvFile(t *testing.T) {
	exam := &EnvFile{Vars: []string{"VAR1", "VAR2"}}

	// Test when environment variables are not set
	result, err := exam.Examinate()
	assert.False(t, result)
	assert.Error(t, err)

	// Test when environment variables are not valid files
	t.Setenv("VAR1", "/invalid/path")
	t.Setenv("VAR2", "/etc/hosts")
	result, err = exam.Examinate()
	assert.False(t, result)
	assert.Error(t, err)

	// Test when environment variables are valid files
	t.Setenv("VAR1", "/etc/hosts")
	result, err = exam.Examinate()
	assert.True(t, result)
	assert.NoError(t, err)
}

func TestEnvDir(t *testing.T) {
	exam := &EnvDir{Vars: []string{"VAR1", "VAR2"}}

	// Test when environment variables are not set
	result, err := exam.Examinate()
	assert.False(t, result)
	assert.Error(t, err)

	// Test when environment variables are not valid directories
	t.Setenv("VAR1", "/invalid/path")
	t.Setenv("VAR2", "/etc")
	result, err = exam.Examinate()
	assert.False(t, result)
	assert.Error(t, err)

	// Test when environment variables are valid directories
	t.Setenv("VAR1", "/etc")
	result, err = exam.Examinate()
	assert.True(t, result)
	assert.NoError(t, err)
}

func TestEnvIpv4Addr(t *testing.T) {
	exam := &EnvIpv4Addr{Vars: []string{"VAR1", "VAR2"}}

	// Test when environment variables are not set
	result, err := exam.Examinate()
	assert.False(t, result)
	assert.Error(t, err)

	// Test when environment variables are not valid IPv4 addresses
	t.Setenv("VAR1", "invalid")
	t.Setenv("VAR2", "192.168.1.1")
	result, err = exam.Examinate()
	assert.False(t, result)
	assert.Error(t, err)

	// Test when environment variables are valid IPv4 addresses
	t.Setenv("VAR1", "10.0.0.1")
	result, err = exam.Examinate()
	assert.True(t, result)
	assert.NoError(t, err)
}

func TestEnvIpv6Addr(t *testing.T) {
	exam := &EnvIpv6Addr{Vars: []string{"VAR1", "VAR2"}}

	// Test when environment variables are not set
	result, err := exam.Examinate()
	assert.False(t, result)
	assert.Error(t, err)

	// Test when environment variables are not valid IPv6 addresses
	t.Setenv("VAR1", "invalid")
	t.Setenv("VAR2", "2001:0db8:85a3:0000:0000:8a2e:0370:7334")
	result, err = exam.Examinate()
	assert.False(t, result)
	assert.Error(t, err)

	// Test when environment variables are valid IPv6 addresses
	t.Setenv("VAR1", "fe80::1ff:fe23:4567:890a")
	result, err = exam.Examinate()
	assert.True(t, result)
	assert.NoError(t, err)
}

func TestEnvIpAddr(t *testing.T) {
	exam := &EnvIpAddr{Vars: []string{"VAR1", "VAR2"}}

	// Test when environment variables are not set
	result, err := exam.Examinate()
	assert.False(t, result)
	assert.Error(t, err)

	// Test when environment variables are not valid IP addresses
	t.Setenv("VAR1", "invalid")
	t.Setenv("VAR2", "192.168.1.1")
	result, err = exam.Examinate()
	assert.False(t, result)
	assert.Error(t, err)

	// Test when environment variables are valid IP addresses
	t.Setenv("VAR1", "2001:0db8:85a3:0000:0000:8a2e:0370:7334")
	result, err = exam.Examinate()
	assert.True(t, result)
	assert.NoError(t, err)
}

func TestEnvHostname(t *testing.T) {
	exam := &EnvHostname{Vars: []string{"VAR1", "VAR2"}, Protocol: "http"}

	// Test when environment variables are not set
	result, err := exam.Examinate()
	assert.False(t, result)
	assert.Error(t, err)

	// Test when environment variables are not valid hostnames
	t.Setenv("VAR1", "invalid")
	t.Setenv("VAR2", "http://example.com")
	result, err = exam.Examinate()
	assert.False(t, result)
	assert.Error(t, err)

	// Test when environment variables are valid hostnames
	t.Setenv("VAR1", "http://example.com")
	result, err = exam.Examinate()
	assert.True(t, result)
	assert.NoError(t, err)

	exam2 := &EnvHostname{Vars: []string{"VAR1", "VAR2"}, Protocol: ""}
	t.Setenv("VAR1", "example.com")
	t.Setenv("VAR2", "tcp://example.com")
	result, err = exam2.Examinate()
	assert.True(t, result)
	assert.NoError(t, err)
	t.Setenv("VAR1", "\n")
	result, err = exam2.Examinate()
	assert.False(t, result)
	assert.Error(t, err)
}

func TestEnvIsSetParse(t *testing.T) {
	exam := &EnvIsSet{}

	// Test invalid type
	_, err := exam.Parse(config.Exam{Type: "invalid"})
	assert.Error(t, err)

	// Test vars not set
	_, err = exam.Parse(config.Exam{Type: "env.is-set"})
	assert.Error(t, err)

	// Test valid config
	parsed, err := exam.Parse(config.Exam{Type: "env.is-set", Vars: []string{"VAR1"}})
	assert.NoError(t, err)
	assert.Equal(t, &EnvIsSet{Vars: []string{"VAR1"}}, parsed)
}

func TestEnvIsSetNotEmptyParse(t *testing.T) {
	exam := &EnvIsSetNotEmpty{}

	// Test invalid type
	_, err := exam.Parse(config.Exam{Type: "invalid"})
	assert.Error(t, err)

	// Test vars not set
	_, err = exam.Parse(config.Exam{Type: "env.not-empty"})
	assert.Error(t, err)

	// Test valid config
	parsed, err := exam.Parse(config.Exam{Type: "env.not-empty", Vars: []string{"VAR1"}})
	assert.NoError(t, err)
	assert.Equal(t, &EnvIsSetNotEmpty{Vars: []string{"VAR1"}}, parsed)
}

func TestEnvRegexParse(t *testing.T) {
	exam := &EnvRegex{}

	// Test invalid type
	_, err := exam.Parse(config.Exam{Type: "invalid"})
	assert.Error(t, err)

	// Test vars not set
	_, err = exam.Parse(config.Exam{Type: "env.regex"})
	assert.Error(t, err)

	// Test regex not set
	_, err = exam.Parse(config.Exam{Type: "env.regex", Vars: []string{"VAR1"}})
	assert.Error(t, err)

	// Test invalid regex
	_, err = exam.Parse(config.Exam{Type: "env.regex", Vars: []string{"VAR1"}, Regex: "["})
	assert.Error(t, err)

	// Test valid config
	parsed, err := exam.Parse(config.Exam{Type: "env.regex", Vars: []string{"VAR1"}, Regex: ".*"})
	assert.NoError(t, err)
	assert.NotNil(t, parsed)
}

func TestEnvOptionParse(t *testing.T) {
	exam := &EnvOption{}

	// Test invalid type
	_, err := exam.Parse(config.Exam{Type: "invalid"})
	assert.Error(t, err)

	// Test vars not set
	_, err = exam.Parse(config.Exam{Type: "env.options"})
	assert.Error(t, err)

	// Test options not set
	_, err = exam.Parse(config.Exam{Type: "env.options", Vars: []string{"VAR1"}})
	assert.Error(t, err)

	// Test valid config
	parsed, err := exam.Parse(config.Exam{Type: "env.options", Vars: []string{"VAR1"}, Options: []string{"option1"}})
	assert.NoError(t, err)
	assert.NotNil(t, parsed)
}

func TestEnvIntegerParse(t *testing.T) {
	exam := &EnvInteger{}

	// Test invalid type
	_, err := exam.Parse(config.Exam{Type: "invalid"})
	assert.Error(t, err)

	// Test vars not set
	_, err = exam.Parse(config.Exam{Type: "env.int"})
	assert.Error(t, err)

	// Test valid config
	parsed, err := exam.Parse(config.Exam{Type: "env.int", Vars: []string{"VAR1"}})
	assert.NoError(t, err)
	assert.Equal(t, &EnvInteger{Vars: []string{"VAR1"}}, parsed)
}

func TestEnvIntegerRangeParse(t *testing.T) {
	exam := &EnvIntegerRange{}

	// Test invalid type
	_, err := exam.Parse(config.Exam{Type: "invalid"})
	assert.Error(t, err)

	// Test vars not set
	_, err = exam.Parse(config.Exam{Type: "env.int-range"})
	assert.Error(t, err)

	// Test min not an integer
	_, err = exam.Parse(config.Exam{Type: "env.int-range", Vars: []string{"VAR1"}, Min: "min", Max: 10})
	assert.Error(t, err)

	// Test max not an integer
	_, err = exam.Parse(config.Exam{Type: "env.int-range", Vars: []string{"VAR1"}, Min: 0, Max: "max"})
	assert.Error(t, err)

	// Test valid config
	parsed, err := exam.Parse(config.Exam{Type: "env.int-range", Vars: []string{"VAR1"}, Min: 0, Max: 10})
	assert.NoError(t, err)
	assert.Equal(t, &EnvIntegerRange{Vars: []string{"VAR1"}, Min: 0, Max: 10}, parsed)
}

func TestEnvFloatParse(t *testing.T) {
	exam := &EnvFloat{}

	// Test invalid type
	_, err := exam.Parse(config.Exam{Type: "invalid"})
	assert.Error(t, err)

	// Test vars not set
	_, err = exam.Parse(config.Exam{Type: "env.float"})
	assert.Error(t, err)

	// Test valid config
	parsed, err := exam.Parse(config.Exam{Type: "env.float", Vars: []string{"VAR1"}})
	assert.NoError(t, err)
	assert.Equal(t, &EnvFloat{Vars: []string{"VAR1"}}, parsed)
}

func TestEnvFloatRangeParse(t *testing.T) {
	exam := &EnvFloatRange{}

	// Test invalid type
	_, err := exam.Parse(config.Exam{Type: "invalid"})
	assert.Error(t, err)

	// Test vars not set
	_, err = exam.Parse(config.Exam{Type: "env.float-range"})
	assert.Error(t, err)

	// Test min not a float
	_, err = exam.Parse(config.Exam{Type: "env.float-range", Vars: []string{"VAR1"}, Min: "min", Max: 10.0})
	assert.Error(t, err)

	// Test max not a float
	_, err = exam.Parse(config.Exam{Type: "env.float-range", Vars: []string{"VAR1"}, Min: 0.0, Max: "max"})
	assert.Error(t, err)

	// Test valid config
	parsed, err := exam.Parse(config.Exam{Type: "env.float-range", Vars: []string{"VAR1"}, Min: 0.0, Max: 10.0})
	assert.NoError(t, err)
	assert.Equal(t, &EnvFloatRange{Vars: []string{"VAR1"}, Min: 0.0, Max: 10.0}, parsed)
}

func TestEnvFileParse(t *testing.T) {
	exam := &EnvFile{}

	// Test invalid type
	_, err := exam.Parse(config.Exam{Type: "invalid"})
	assert.Error(t, err)

	// Test vars not set
	_, err = exam.Parse(config.Exam{Type: "env.file"})
	assert.Error(t, err)

	// Test valid config
	parsed, err := exam.Parse(config.Exam{Type: "env.file", Vars: []string{"VAR1"}})
	assert.NoError(t, err)
	assert.Equal(t, &EnvFile{Vars: []string{"VAR1"}}, parsed)
}

func TestEnvDirParse(t *testing.T) {
	exam := &EnvDir{}

	// Test invalid type
	_, err := exam.Parse(config.Exam{Type: "invalid"})
	assert.Error(t, err)

	// Test vars not set
	_, err = exam.Parse(config.Exam{Type: "env.dir"})
	assert.Error(t, err)

	// Test valid config
	parsed, err := exam.Parse(config.Exam{Type: "env.dir", Vars: []string{"VAR1"}})
	assert.NoError(t, err)
	assert.Equal(t, &EnvDir{Vars: []string{"VAR1"}}, parsed)
}

func TestEnvIpv4AddrParse(t *testing.T) {
	exam := &EnvIpv4Addr{}

	// Test invalid type
	_, err := exam.Parse(config.Exam{Type: "invalid"})
	assert.Error(t, err)

	// Test vars not set
	_, err = exam.Parse(config.Exam{Type: "env.ipv4"})
	assert.Error(t, err)

	// Test valid config
	parsed, err := exam.Parse(config.Exam{Type: "env.ipv4", Vars: []string{"VAR1"}})
	assert.NoError(t, err)
	assert.Equal(t, &EnvIpv4Addr{Vars: []string{"VAR1"}}, parsed)
}

func TestEnvIpv6AddrParse(t *testing.T) {
	exam := &EnvIpv6Addr{}

	// Test invalid type
	_, err := exam.Parse(config.Exam{Type: "invalid"})
	assert.Error(t, err)

	// Test vars not set
	_, err = exam.Parse(config.Exam{Type: "env.ipv6"})
	assert.Error(t, err)

	// Test valid config
	parsed, err := exam.Parse(config.Exam{Type: "env.ipv6", Vars: []string{"VAR1"}})
	assert.NoError(t, err)
	assert.Equal(t, &EnvIpv6Addr{Vars: []string{"VAR1"}}, parsed)
}

func TestEnvIpAddrParse(t *testing.T) {
	exam := &EnvIpAddr{}

	// Test invalid type
	_, err := exam.Parse(config.Exam{Type: "invalid"})
	assert.Error(t, err)

	// Test vars not set
	_, err = exam.Parse(config.Exam{Type: "env.ip"})
	assert.Error(t, err)

	// Test valid config
	parsed, err := exam.Parse(config.Exam{Type: "env.ip", Vars: []string{"VAR1"}})
	assert.NoError(t, err)
	assert.Equal(t, &EnvIpAddr{Vars: []string{"VAR1"}}, parsed)
}

func TestEnvHostnameParse(t *testing.T) {
	exam := &EnvHostname{}

	// Test invalid type
	_, err := exam.Parse(config.Exam{Type: "invalid"})
	assert.Error(t, err)

	// Test vars not set
	_, err = exam.Parse(config.Exam{Type: "env.hostname"})
	assert.Error(t, err)

	// Test valid config
	parsed, err := exam.Parse(config.Exam{Type: "env.hostname", Vars: []string{"VAR1"}, Protocol: "http"})
	assert.NoError(t, err)
	assert.Equal(t, &EnvHostname{Vars: []string{"VAR1"}, Protocol: "http"}, parsed)
}
