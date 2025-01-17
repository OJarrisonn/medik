package env

import (
	"regexp"
	"testing"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/medik"
	"github.com/stretchr/testify/assert"
)

func TestEnvParsersCollection(t *testing.T) {
	// Inexistent parser
	parser, ok := GetParser("inexistent")

	assert.Nil(t, parser)
	assert.False(t, ok)

	// Existing parsers
	parser, ok = GetParser("env.is-set")

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
		(&IsSet{}).Type(),
		(&NotEmpty{}).Type(),
		(&Regex{}).Type(),
		(&Option{}).Type(),
		(&Int{}).Type(),
		(&IntRange{}).Type(),
		(&Float{}).Type(),
		(&FloatRange{}).Type(),
		(&File{}).Type(),
		(&Dir{}).Type(),
		(&Ipv4{}).Type(),
		(&Ipv6{}).Type(),
		(&Ip{}).Type(),
		(&Hostname{}).Type(),
	}

	assert.ElementsMatch(t, known, registered)
}

func TestEnvWrongParser(t *testing.T) {
	parse, ok := GetParser("env.is-set")

	assert.True(t, ok)

	_, err := parse(config.Exam{Type: "invalid"})

	assert.NotNil(t, err)
	assert.Equal(t, "wrong exam parser: using env.is-set parser for a invalid exam", err.Error())
}

func TestEnvIsSetNotEmpty(t *testing.T) {
	exam := &NotEmpty{Vars: []string{"VAR1", "VAR2"}, Level: medik.ERROR}

	// Test when environment variables are not set
	report := exam.Examinate()
	ok, header, body := report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables are set to empty
	t.Setenv("VAR1", "")
	t.Setenv("VAR2", " ")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables are set to non-empty values
	t.Setenv("VAR1", "value1")
	t.Setenv("VAR2", "value2")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.OK, ok)
	assert.NotEmpty(t, header)
	assert.Empty(t, body)
}

func TestEnvRegex(t *testing.T) {
	regex, _ := regexp.Compile(`^value\d$`)
	exam := &Regex{Vars: []string{"VAR1", "VAR2"}, Regex: regex, Level: medik.ERROR}

	// Test when environment variables are not set
	report := exam.Examinate()
	ok, header, body := report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables do not match regex
	t.Setenv("VAR1", "invalid")
	t.Setenv("VAR2", "value2")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables match regex
	t.Setenv("VAR1", "value1")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.OK, ok)
	assert.NotEmpty(t, header)
	assert.Empty(t, body)
}

func TestEnvOption(t *testing.T) {
	options := map[string]bool{"option1": true, "option2": true}
	exam := &Option{Vars: []string{"VAR1", "VAR2"}, Options: options, Level: medik.ERROR}

	// Test when environment variables are not set
	report := exam.Examinate()
	ok, header, body := report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables do not match options
	t.Setenv("VAR1", "invalid")
	t.Setenv("VAR2", "option2")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables match options
	t.Setenv("VAR1", "option1")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.OK, ok)
	assert.NotEmpty(t, header)
	assert.Empty(t, body)
}

func TestEnvInteger(t *testing.T) {
	exam := &Int{Vars: []string{"VAR1", "VAR2"}, Level: medik.ERROR}

	// Test when environment variables are not set
	report := exam.Examinate()
	ok, header, body := report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables are not integers
	t.Setenv("VAR1", "invalid")
	t.Setenv("VAR2", "123")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables are integers
	t.Setenv("VAR1", "456")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.OK, ok)
	assert.NotEmpty(t, header)
	assert.Empty(t, body)
}

func TestEnvIntegerRange(t *testing.T) {
	exam := &IntRange{Vars: []string{"VAR1", "VAR2"}, Min: 10, Max: 100, Level: medik.ERROR}

	// Test when environment variables are not set
	report := exam.Examinate()
	ok, header, body := report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables are not integers
	t.Setenv("VAR1", "invalid")
	t.Setenv("VAR2", "50")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables are out of range
	t.Setenv("VAR1", "5")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables are within range
	t.Setenv("VAR1", "20")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.OK, ok)
	assert.NotEmpty(t, header)
	assert.Empty(t, body)
}

func TestEnvFloat(t *testing.T) {
	exam := &Float{Vars: []string{"VAR1", "VAR2"}, Level: medik.ERROR}

	// Test when environment variables are not set
	report := exam.Examinate()
	ok, header, body := report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables are not floats
	t.Setenv("VAR1", "invalid")
	t.Setenv("VAR2", "123.45")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables are floats
	t.Setenv("VAR1", "456.78")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.OK, ok)
	assert.NotEmpty(t, header)
	assert.Empty(t, body)
}

func TestEnvFloatRange(t *testing.T) {
	exam := &FloatRange{Vars: []string{"VAR1", "VAR2"}, Min: 10.5, Max: 100.5, Level: medik.ERROR}

	// Test when environment variables are not set
	report := exam.Examinate()
	ok, header, body := report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables are not floats
	t.Setenv("VAR1", "invalid")
	t.Setenv("VAR2", "50.5")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables are out of range
	t.Setenv("VAR1", "5.5")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables are within range
	t.Setenv("VAR1", "20.5")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.OK, ok)
	assert.NotEmpty(t, header)
	assert.Empty(t, body)
}

func TestEnvFileExists(t *testing.T) {
	exam := &File{Vars: []string{"VAR1", "VAR2"}, Exists: true, Level: medik.ERROR}

	// Test when environment variables are not set
	report := exam.Examinate()
	ok, header, body := report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables are not valid files
	t.Setenv("VAR1", "/invalid/path")
	t.Setenv("VAR2", "/etc/hosts")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables are valid files
	t.Setenv("VAR1", "/etc/hosts")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.OK, ok)
	assert.NotEmpty(t, header)
	assert.Empty(t, body)
}

func TestEnvFileNotExists(t *testing.T) {
	exam := &File{Vars: []string{"VAR1", "VAR2"}, Exists: false, Level: medik.ERROR}

	// Test when environment variables are not set
	report := exam.Examinate()
	ok, header, body := report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables are valid files
	t.Setenv("VAR1", "/etc/hosts")
	t.Setenv("VAR2", "/invalid/path")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables are not valid files
	t.Setenv("VAR1", "/invalid/path")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.OK, ok)
	assert.NotEmpty(t, header)
	assert.Empty(t, body)
}

func TestEnvIpv4Addr(t *testing.T) {
	exam := &Ipv4{Vars: []string{"VAR1", "VAR2"}, Level: medik.ERROR}
	// Test when environment variables are not set
	report := exam.Examinate()
	ok, header, body := report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables are not valid IPv4 addresses
	t.Setenv("VAR1", "invalid")
	t.Setenv("VAR2", "192.168.1.1")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables are valid IPv4 addresses
	t.Setenv("VAR1", "10.0.0.1")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.OK, ok)
	assert.NotEmpty(t, header)
	assert.Empty(t, body)
}

func TestEnvIpv6Addr(t *testing.T) {
	exam := &Ipv6{Vars: []string{"VAR1", "VAR2"}, Level: medik.ERROR}

	// Test when environment variables are not set
	report := exam.Examinate()
	ok, header, body := report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables are not valid IPv6 addresses
	t.Setenv("VAR1", "invalid")
	t.Setenv("VAR2", "2001:0db8:85a3:0000:0000:8a2e:0370:7334")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables are valid IPv6 addresses
	t.Setenv("VAR1", "fe80::1ff:fe23:4567:890a")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.OK, ok)
	assert.NotEmpty(t, header)
	assert.Empty(t, body)
}

func TestEnvIpAddr(t *testing.T) {
	exam := &Ip{Vars: []string{"VAR1", "VAR2"}, Level: medik.ERROR}

	// Test when environment variables are not set
	report := exam.Examinate()
	ok, header, body := report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables are not valid IP addresses
	t.Setenv("VAR1", "invalid")
	t.Setenv("VAR2", "192.168.1.1")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables are valid IP addresses
	t.Setenv("VAR1", "fe80::1ff:fe23:4567:890a")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.OK, ok)
	assert.NotEmpty(t, header)
	assert.Empty(t, body)
}

func TestEnvHostname(t *testing.T) {
	exam := &Hostname{Vars: []string{"VAR1", "VAR2"}, Protocol: "http", Level: medik.ERROR}

	// Test when environment variables are not set
	report := exam.Examinate()
	ok, header, body := report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables are not valid hostnames
	t.Setenv("VAR1", "invalid")
	t.Setenv("VAR2", "http://example.com")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)

	// Test when environment variables are valid hostnames
	t.Setenv("VAR1", "http://example.com")
	report = exam.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.OK, ok)
	assert.NotEmpty(t, header)
	assert.Empty(t, body)

	exam2 := &Hostname{Vars: []string{"VAR1", "VAR2"}, Protocol: "", Level: medik.ERROR}
	t.Setenv("VAR1", "example.com")
	t.Setenv("VAR2", "tcp://example.com")
	report = exam2.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.OK, ok)
	assert.NotEmpty(t, header)
	assert.Empty(t, body)
	t.Setenv("VAR1", "\n")
	report = exam2.Examinate()
	ok, header, body = report.Format(medik.WARNING)
	assert.Equal(t, medik.ERROR, ok)
	assert.NotEmpty(t, header)
	assert.NotEmpty(t, body)
}

func TestEnvIsSetParse(t *testing.T) {
	exam := &IsSet{Level: medik.ERROR}

	// Test invalid type
	_, err := exam.Parse(config.Exam{Type: "invalid"})
	assert.NotNil(t, err)

	// Test vars not set
	_, err = exam.Parse(config.Exam{Type: "env.is-set"})
	assert.NotNil(t, err)

	// Test valid config
	parsed, err := exam.Parse(config.Exam{Type: "env.is-set", Vars: []string{"VAR1"}})
	assert.Nil(t, err)
	assert.Equal(t, &IsSet{Vars: []string{"VAR1"}, Level: medik.ERROR}, parsed)
}

func TestEnvIsSetNotEmptyParse(t *testing.T) {
	exam := &NotEmpty{Level: medik.ERROR}

	// Test invalid type
	_, err := exam.Parse(config.Exam{Type: "invalid"})
	assert.NotNil(t, err)

	// Test vars not set
	_, err = exam.Parse(config.Exam{Type: "env.not-empty"})
	assert.NotNil(t, err)

	// Test valid config
	parsed, err := exam.Parse(config.Exam{Type: "env.not-empty", Vars: []string{"VAR1"}})
	assert.Nil(t, err)
	assert.Equal(t, &NotEmpty{Vars: []string{"VAR1"}, Level: medik.ERROR}, parsed)
}

func TestEnvRegexParse(t *testing.T) {
	exam := &Regex{Level: medik.ERROR}

	// Test invalid type
	_, err := exam.Parse(config.Exam{Type: "invalid"})
	assert.NotNil(t, err)

	// Test vars not set
	_, err = exam.Parse(config.Exam{Type: "env.regex"})
	assert.NotNil(t, err)

	// Test regex not set
	_, err = exam.Parse(config.Exam{Type: "env.regex", Vars: []string{"VAR1"}})
	assert.NotNil(t, err)

	// Test invalid regex
	_, err = exam.Parse(config.Exam{Type: "env.regex", Vars: []string{"VAR1"}, Regex: "["})
	assert.NotNil(t, err)

	// Test valid config
	parsed, err := exam.Parse(config.Exam{Type: "env.regex", Vars: []string{"VAR1"}, Regex: ".*"})
	assert.Nil(t, err)
	assert.NotNil(t, parsed)
}

func TestEnvOptionParse(t *testing.T) {
	exam := &Option{Level: medik.ERROR}

	// Test invalid type
	_, err := exam.Parse(config.Exam{Type: "invalid"})
	assert.NotNil(t, err)

	// Test vars not set
	_, err = exam.Parse(config.Exam{Type: "env.options"})
	assert.NotNil(t, err)

	// Test options not set
	_, err = exam.Parse(config.Exam{Type: "env.options", Vars: []string{"VAR1"}})
	assert.NotNil(t, err)

	// Test valid config
	parsed, err := exam.Parse(config.Exam{Type: "env.options", Vars: []string{"VAR1"}, Options: []string{"option1"}})
	assert.Nil(t, err)
	assert.NotNil(t, parsed)
}

func TestEnvIntegerParse(t *testing.T) {
	exam := &Int{Level: medik.ERROR}

	// Test invalid type
	_, err := exam.Parse(config.Exam{Type: "invalid"})
	assert.NotNil(t, err)

	// Test vars not set
	_, err = exam.Parse(config.Exam{Type: "env.int"})
	assert.NotNil(t, err)

	// Test valid config
	parsed, err := exam.Parse(config.Exam{Type: "env.int", Vars: []string{"VAR1"}})
	assert.Nil(t, err)
	assert.Equal(t, &Int{Vars: []string{"VAR1"}, Level: medik.ERROR}, parsed)
}

func TestEnvIntegerRangeParse(t *testing.T) {
	exam := &IntRange{Level: medik.ERROR}

	// Test invalid type
	_, err := exam.Parse(config.Exam{Type: "invalid"})
	assert.NotNil(t, err)

	// Test vars not set
	_, err = exam.Parse(config.Exam{Type: "env.int-range"})
	assert.NotNil(t, err)

	// Test min not an integer
	_, err = exam.Parse(config.Exam{Type: "env.int-range", Vars: []string{"VAR1"}, Min: "min", Max: 10})
	assert.NotNil(t, err)

	// Test max not an integer
	_, err = exam.Parse(config.Exam{Type: "env.int-range", Vars: []string{"VAR1"}, Min: 0, Max: "max"})
	assert.NotNil(t, err)

	// Test valid config
	parsed, err := exam.Parse(config.Exam{Type: "env.int-range", Vars: []string{"VAR1"}, Min: 0, Max: 10})
	assert.Nil(t, err)
	assert.Equal(t, &IntRange{Vars: []string{"VAR1"}, Min: 0, Max: 10, Level: medik.ERROR}, parsed)
}

func TestEnvFloatParse(t *testing.T) {
	exam := &Float{Level: medik.ERROR}

	// Test invalid type
	_, err := exam.Parse(config.Exam{Type: "invalid"})
	assert.NotNil(t, err)

	// Test vars not set
	_, err = exam.Parse(config.Exam{Type: "env.float"})
	assert.NotNil(t, err)

	// Test valid config
	parsed, err := exam.Parse(config.Exam{Type: "env.float", Vars: []string{"VAR1"}})
	assert.Nil(t, err)
	assert.Equal(t, &Float{Vars: []string{"VAR1"}, Level: medik.ERROR}, parsed)
}

func TestEnvFloatRangeParse(t *testing.T) {
	exam := &FloatRange{Level: medik.ERROR}

	// Test invalid type
	_, err := exam.Parse(config.Exam{Type: "invalid"})
	assert.NotNil(t, err)

	// Test vars not set
	_, err = exam.Parse(config.Exam{Type: "env.float-range"})
	assert.NotNil(t, err)

	// Test min not a float
	_, err = exam.Parse(config.Exam{Type: "env.float-range", Vars: []string{"VAR1"}, Min: "min", Max: 10.0})
	assert.NotNil(t, err)

	// Test max not a float
	_, err = exam.Parse(config.Exam{Type: "env.float-range", Vars: []string{"VAR1"}, Min: 0.0, Max: "max"})
	assert.NotNil(t, err)

	// Test valid config
	parsed, err := exam.Parse(config.Exam{Type: "env.float-range", Vars: []string{"VAR1"}, Min: 0.0, Max: 10.0})
	assert.Nil(t, err)
	assert.Equal(t, &FloatRange{Vars: []string{"VAR1"}, Min: 0.0, Max: 10.0, Level: medik.ERROR}, parsed)
}

func TestEnvFileParse(t *testing.T) {
	exam := &File{Level: medik.ERROR}

	// Test invalid type
	_, err := exam.Parse(config.Exam{Type: "invalid"})
	assert.NotNil(t, err)

	// Test vars not set
	_, err = exam.Parse(config.Exam{Type: "env.file"})
	assert.NotNil(t, err)

	// Test valid config
	parsed, err := exam.Parse(config.Exam{Type: "env.file", Vars: []string{"VAR1"}})
	assert.Nil(t, err)
	assert.Equal(t, &File{Vars: []string{"VAR1"}, Level: medik.ERROR}, parsed)
}

func TestEnvDirParse(t *testing.T) {
	exam := &Dir{Level: medik.ERROR}

	// Test invalid type
	_, err := exam.Parse(config.Exam{Type: "invalid"})
	assert.NotNil(t, err)

	// Test vars not set
	_, err = exam.Parse(config.Exam{Type: "env.dir"})
	assert.NotNil(t, err)

	// Test valid config
	parsed, err := exam.Parse(config.Exam{Type: "env.dir", Vars: []string{"VAR1"}})
	assert.Nil(t, err)
	assert.Equal(t, &Dir{Vars: []string{"VAR1"}, Level: medik.ERROR}, parsed)
}

func TestEnvIpv4AddrParse(t *testing.T) {
	exam := &Ipv4{Level: medik.ERROR}

	// Test invalid type
	_, err := exam.Parse(config.Exam{Type: "invalid"})
	assert.NotNil(t, err)

	// Test vars not set
	_, err = exam.Parse(config.Exam{Type: "env.ipv4"})
	assert.NotNil(t, err)

	// Test valid config
	parsed, err := exam.Parse(config.Exam{Type: "env.ipv4", Vars: []string{"VAR1"}})
	assert.Nil(t, err)
	assert.Equal(t, &Ipv4{Vars: []string{"VAR1"}, Level: medik.ERROR}, parsed)
}

func TestEnvIpv6AddrParse(t *testing.T) {
	exam := &Ipv6{Level: medik.ERROR}

	// Test invalid type
	_, err := exam.Parse(config.Exam{Type: "invalid"})
	assert.NotNil(t, err)

	// Test vars not set
	_, err = exam.Parse(config.Exam{Type: "env.ipv6"})
	assert.NotNil(t, err)

	// Test valid config
	parsed, err := exam.Parse(config.Exam{Type: "env.ipv6", Vars: []string{"VAR1"}})
	assert.Nil(t, err)
	assert.Equal(t, &Ipv6{Vars: []string{"VAR1"}, Level: medik.ERROR}, parsed)
}

func TestEnvIpAddrParse(t *testing.T) {
	exam := &Ip{Level: medik.ERROR}

	// Test invalid type
	_, err := exam.Parse(config.Exam{Type: "invalid"})
	assert.NotNil(t, err)

	// Test vars not set
	_, err = exam.Parse(config.Exam{Type: "env.ip"})
	assert.NotNil(t, err)

	// Test valid config
	parsed, err := exam.Parse(config.Exam{Type: "env.ip", Vars: []string{"VAR1"}})
	assert.Nil(t, err)
	assert.Equal(t, &Ip{Vars: []string{"VAR1"}, Level: medik.ERROR}, parsed)
}

func TestEnvHostnameParse(t *testing.T) {
	exam := &Hostname{Level: medik.ERROR}

	// Test invalid type
	_, err := exam.Parse(config.Exam{Type: "invalid"})
	assert.NotNil(t, err)

	// Test vars not set
	_, err = exam.Parse(config.Exam{Type: "env.hostname"})
	assert.NotNil(t, err)

	// Test valid config
	parsed, err := exam.Parse(config.Exam{Type: "env.hostname", Vars: []string{"VAR1"}, Protocol: "http"})
	assert.Nil(t, err)
	assert.Equal(t, &Hostname{Vars: []string{"VAR1"}, Protocol: "http", Level: medik.ERROR}, parsed)
}
