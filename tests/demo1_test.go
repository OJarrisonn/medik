package tests

import (
	"testing"

	"github.com/OJarrisonn/medik/pkg/cmd"
)

func TestRunDemo1(t *testing.T) {
	cmd.ConfigFile = "../samples/medik.demo1.yaml"
	cmd.EnvFile = "../samples/.env.demo1"

	cmd.Execute(nil)
}
