package cmd

import (
	"os"
	"testing"
)

func TestUseEmptyEnv(t *testing.T) {
	_, err := useEnv("")
	if err != nil {
		t.Errorf("useEnv() failed: %v", err)
	}
}

func TestUseEnv(t *testing.T) {
	_, err := useEnv("FOO=bar")
	if err != nil {
		t.Errorf("useEnv() failed: %v", err)
	}

	if os.Getenv("FOO") != "bar" {
		t.Errorf("useEnv() failed: %v", err)
	}

	os.Unsetenv("FOO")
}

func TestUseEnvWithComment(t *testing.T) {
	_, err := useEnv("# FOO=bar")
	if err != nil {
		t.Errorf("useEnv() failed: % v", err)
	}

	if os.Getenv("FOO") != "" {
		t.Errorf("useEnv() failed: FOO is set")
	}

	os.Unsetenv("FOO")
}

func TestUseEnvWithMultilines(t *testing.T) {
	_, err := useEnv("FOO=bar\nBAR=baz")
	if err != nil {
		t.Errorf("useEnv() failed: %v", err)
	}

	if os.Getenv("FOO") != "bar" {
		t.Errorf("useEnv() failed: %v", err)
	}

	if os.Getenv("BAR") != "baz" {
		t.Errorf("useEnv() failed: %v", err)
	}

	os.Unsetenv("FOO")
	os.Unsetenv("BAR")
}

func TestUseEnvWithQuotes(t *testing.T) {
	_, err := useEnv("FOO=\"bar\"")
	if err != nil {
		t.Errorf("useEnv() failed: %v", err)
	}

	if val := os.Getenv("FOO"); val != "bar" {
		t.Errorf("useEnv() failed: FOO = '%v'", val)
	}

	os.Unsetenv("FOO")
}

func TestUseEnvWithQuotesAndComment(t *testing.T) {
	_, err := useEnv("FOO=\"bar\" # baz")
	if err != nil {
		t.Errorf("useEnv() failed: %v", err)
	}

	if val := os.Getenv("FOO"); val != "bar" {
		t.Errorf("useEnv() failed: FOO = '%v'", val)
	}

	os.Unsetenv("FOO")
}
