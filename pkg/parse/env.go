package parse

import (
	"fmt"
	"strings"
)

func ParseEnvFile(content string) (map[string]string, error) {
	lines := strings.Split(content, "\n")
	lines = cleanLines(lines)

	env := map[string]string{}

	for _, line := range lines {
		key, value, err := splitKeyValue(line)
		if err != nil {
			return nil, err
		}

		env[key] = value
	}

	return env, nil
}

func cleanLines(lines []string) []string {
	cleaned := []string{}

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if isLineIgnorable(line) {
			continue
		}

		line = trimComment(line)

		if line != "" {
			cleaned = append(cleaned, line)
		}
	}

	return cleaned
}

func isLineIgnorable(line string) bool {
	return line == "" || strings.HasPrefix(line, "#")
}

func trimComment(line string) string {
	// state = ' ': may finish, '\'': single quote started, '"': double quote started, '\\': escape character
	lastState := ' '
	state := ' '

loop:
	for i, c := range line {
		tmp := state
		switch state {
		case ' ':
			if c == '#' {
				line = line[:i]
				break loop
			} else if c == '\'' {
				state = '\''
			} else if c == '"' {
				state = '"'
			}
		case '\\':
			state = lastState
		case '\'':
			switch c {
			case '\'':
				state = ' '
			case '\\':
				state = '\\'
			}
		case '"':
			switch c {
			case '"':
				state = ' '
			case '\\':
				state = '\\'
			}
		}
		lastState = tmp
	}

	return strings.TrimSpace(line)
}

func splitKeyValue(line string) (string, string, error) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid line: %s", line)
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	if strings.HasPrefix(value, "'") || strings.HasPrefix(value, "\"") {
		value = value[1:]
	}

	if strings.HasSuffix(value, "'") || strings.HasSuffix(value, "\"") {
		value = value[:len(value)-1]
	}

	if key == "" {
		return "", "", fmt.Errorf("invalid line: %s", line)
	}

	return key, value, nil
}
