package tests

import (
	"testing"

	"github.com/OJarrisonn/medik/pkg/cmd"
)

func TestRunDemo2(t *testing.T) {
	cmd.ConfigFile = "../samples/medik.demo2.yaml"
	cmd.EnvFile = "../samples/.env.demo2"

	cmd.Execute([]string{"dingle-bell"})
}
