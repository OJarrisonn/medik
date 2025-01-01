package tests

import (
	"testing"

	"github.com/OJarrisonn/medik/pkg/cmd"
	"github.com/OJarrisonn/medik/pkg/medik"
)

func TestRunDemo1(t *testing.T) {
	medik.ConfigFile = "../samples/medik.demo1.yaml"
	medik.EnvFile = "../samples/.env.demo1"

	cmd.Execute(nil)
}
