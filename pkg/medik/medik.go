package medik

import (
	"strings"

	"github.com/fatih/color"
)

const (
	Name        = "medik"
	Version     = "alpha"
	Description = "Medik is a tool for running health checks on a system"
	Author      = "OJarrisonn <jhmtv10@gmail.com>"
)

const (
	DefaultConfigFile = "medik.yaml"
	DefaultEnvFile    = ""
	DefaultVerbose    = false
	DefaultNoColor    = false
)

const (
	OK = iota
	WARNING
	ERROR
)

const (
	MAX_LEVEL     = ERROR
	DEFAULT_LEVEL = "error"
)

var levels = []string{"OK", "WARNING", "ERROR"}

func LogLevel(level int) string {
	return levels[level]
}

func LogLevelFromStr(level string) int {
	level = strings.ToUpper(level)
	for i, l := range levels {
		if l == level {
			return i
		}
	}

	return MAX_LEVEL
}

var (
	ErrorWithBgColor   = color.New(color.BgRed, color.FgBlack)
	WarningWithBgColor = color.New(color.BgYellow, color.FgBlack)
	SuccessWithBgColor = color.New(color.BgGreen, color.FgBlack)
	ErrorColor         = color.New(color.FgRed)
	WarningColor       = color.New(color.FgYellow)
	SuccessColor       = color.New(color.FgGreen)
)

var (
	ConfigFile string
	EnvFile    string
	Verbosity  int
	NoColor    bool
)
