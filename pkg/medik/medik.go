package medik

import "github.com/fatih/color"

const Name = "medik"
const Version = "alpha"
const Description = "Medik is a tool for running health checks on a system"
const Author = "OJarrisonn <jhmtv10@gmail.com>"

const DefaultConfigFile = "medik.yaml"
const DefaultEnvFile = ".env"
const DefaultNoUseEnv = false
const DefaultVerbose = false
const DefaultNoColor = false

const (
	SUCCESS = "OK"
	WARNING = "WARN"
	FAILURE = "ERROR"
)

var ErrorWithBgColor = color.New(color.BgRed, color.FgBlack)
var SuccessWithBgColor = color.New(color.BgGreen, color.FgBlack)
var ErrorColor = color.New(color.FgRed)
var SuccessColor = color.New(color.FgGreen)

var ConfigFile string
var EnvFile string
var NoUseEnv bool
var Verbose bool
var NoColor bool
