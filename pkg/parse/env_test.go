package parse

import (
	"fmt"
	"testing"
)

func TestIsLineIgnorable(t *testing.T) {
	lines := []struct {
		Line      string
		Ignorable bool
	}{
		{"", true},
		{"#", true},
		{"# comment", true},
		{"key=value", false},
		{"key = value", false},
		{" key = value ", false},
		{"key = value # comment", false},
		{"key = # comment", false},
	}

	for _, l := range lines {
		if isLineIgnorable(l.Line) != l.Ignorable {
			t.Errorf("isLineIgnorable() failed: %v", l)
		}
	}
}

func TestTrimComment(t *testing.T) {
	lines := []struct {
		Line   string
		Result string
	}{
		{"key=value", "key=value"},
		{"key='value'", "key='value'"},
		{"key = value", "key = value"},
		{" key = value ", "key = value"},
		{"key = value # comment", "key = value"},
		{"key = # comment", "key ="},
		{"key = \"\\'value\"# comment", "key = \"\\'value\""},
		{"key = '\\'value'# comment", "key = '\\'value'"},
	}

	for _, l := range lines {
		if trimComment(l.Line) != l.Result {
			t.Errorf("trimComment() failed: %v", l)
		}
	}
}

func TestCleanLines(t *testing.T) {
	lines := []struct {
		Lines  []string
		Result []string
	}{
		{[]string{"  ", "# comment", " key=value", ""}, []string{"key=value"}},
		{[]string{"key=value", "# comment", "key=value "}, []string{"key=value", "key=value"}},
		{[]string{"\tkey=value", "# comment", "key=value # comment"}, []string{"key=value", "key=value"}},
		{[]string{"\tkey='\\'value'", "# comment", "key=\"value\" # comment"}, []string{"key='\\'value'", "key=\"value\""}},
	}

	for _, l := range lines {
		if result := cleanLines(l.Lines); len(result) != len(l.Result) {
			t.Errorf("cleanLines() failed: %v :: %v", l, result)
		}
	}
}

func TestSplitKeyValue(t *testing.T) {
	lines := []struct {
		Line   string
		Key    string
		Value  string
		Result error
	}{
		{"key=value", "key", "value", nil},
		{"key = value", "key", "value", nil},
		{" key = value ", "key", "value", nil},
		{"key = value ", "key", "value", nil},
		{"key = 'value' ", "key", "value", nil},
		{"key = \"value\" ", "key", "value", nil},
		{"key = \"value \" ", "key", "value ", nil},
		{"key = ", "key", "", nil},
		{"key", "", "", fmt.Errorf("invalid line: %s", "key")},
		{"=value", "", "", fmt.Errorf("invalid line: %s", "=value")},
	}

	for _, l := range lines {
		key, value, result := splitKeyValue(l.Line)

		if key != l.Key || value != l.Value || result == nil && l.Result != nil || result != nil && l.Result == nil {

			t.Errorf("splitKeyValue() failed: %v :: {%v %v %v}", l, key, value, result)
		}
	}
}

func TestParseEnvFile(t *testing.T) {
	content := `# comment
key=value
key2=value2 # some comment
key3 = value3
# just comment
key4  = 'value4'
`
	env, err := ParseEnvFile(content)
	expected := map[string]string{
		"key":  "value",
		"key2": "value2",
		"key3": "value3",
		"key4": "value4",
	}

	if err != nil {
		t.Errorf("ParseEnvFile() failed: %v", err)
	}

	if len(env) != len(expected) {
		t.Errorf("ParseEnvFile() failed: %v", env)
	}

	for k, v := range env {
		if expected[k] != v {
			t.Errorf("ParseEnvFile() failed: %v: %v, %v", k, v, expected[k])
		}
	}

	content = `=value`

	_, err = ParseEnvFile(content)

	if err == nil {
		t.Errorf("ParseEnvFile() failed: %v", content)
	}
}
