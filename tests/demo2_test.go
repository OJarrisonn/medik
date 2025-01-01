package tests

import (
	"testing"

	"github.com/OJarrisonn/medik/pkg/cmd"
	"github.com/OJarrisonn/medik/pkg/medik"
)

func TestRunDemo2(t *testing.T) {
	medik.ConfigFile = "../samples/medik.demo2.yaml"
	medik.EnvFile = "../samples/.env.demo2"

	cmd.Execute([]string{"dingle-bell"})
}
