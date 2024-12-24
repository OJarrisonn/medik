package config

import "testing"

func TestParseConfigFile(t *testing.T) {
	cfg := `
vitals:
      - type: exists
        vars:
          - DATABASE_URL
          - PORT
`

	_, err := Parse(cfg)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParseConfigFileInvalid(t *testing.T) {
	cfg := `
vitals:
	- tipe: exists
`

	_, err := Parse(cfg)

	if err == nil {
		t.Errorf("expected error, got nil")
	}
}