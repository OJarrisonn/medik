package cmd

import (
	"os"
	"testing"

	"github.com/OJarrisonn/medik/pkg/medik"
)

func TestLoadEnvFileNotSet(t *testing.T) {
	medik.EnvFile = ""

	_, err := loadEnv()
	if err != nil {
		t.Errorf("loadEnv() not accepted empty filename: %v", err)
	}
}

func TestLoadEnvFileInexistent(t *testing.T) {
	medik.EnvFile = "/this/file/is/inexistent.env"

	_, err := loadEnv()
	if err == nil {
		t.Error("loadEnv() accepted an non existent file")
	}
}

func TestLoadEnvFile(t *testing.T) {
	medik.EnvFile = "../../samples/root_test.env"

	_, err := loadEnv()
	if err != nil {
		t.Errorf("loadEnv() failed: %v", err)
	}

	root, ok := os.LookupEnv("ROOT")

	if !ok {
		t.Error("loadEnv() failed: ROOT not set")
	}

	if root != "test" {
		t.Errorf("loadEnv() failed: ROOT = '%v'", root)
	}

	test, ok := os.LookupEnv("TEST")

	if !ok {
		t.Error("loadEnv() failed: TEST not set")
	}

	if test != "root" {
		t.Errorf("loadEnv() failed: TEST = '%v'", test)
	}

	os.Unsetenv("TEST")
	os.Unsetenv("ROOT")
}

func TestUseEmptyEnv(t *testing.T) {
	_, err := useEnv("")
	if err != nil {
		t.Errorf("useEnv() failed: %v", err)
	}
}

func TestUseEnvWithInvalid(t *testing.T) {
	ok, err := useEnv("FOO")

	if err == nil || ok {
		t.Errorf("useEnv() accepted invalid env file: %v", err)
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
