package cmd

import (
	"fmt"
	"os"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/medik"
	"github.com/OJarrisonn/medik/pkg/parse"
	"github.com/OJarrisonn/medik/pkg/runner"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     medik.Name,
	Version: medik.Version,
	Short:   medik.Description,
	Args:    cobra.ArbitraryArgs,
	Run:     run,
}

var ConfigFile string
var EnvFile string
var NoUseEnv bool

func init() {
	rootCmd.PersistentFlags().StringVarP(&ConfigFile, "config", "c", medik.DefaultConfigFile, "Config file to use")
	rootCmd.PersistentFlags().StringVarP(&EnvFile, "env", "e", medik.DefaultEnvFile, "Env file to use")
	rootCmd.PersistentFlags().BoolVar(&NoUseEnv, "no-env", medik.DefaultNoUseEnv, "Won't use an env file")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	cfg, err := loadConfig()

	if err != nil {
		fmt.Printf("Error loading config: %s\n", err)
		os.Exit(1)
	}

	success, errs := runner.Run(cfg, args)

	for _, e := range errs {
		fmt.Println(e)
	}

	if success {
		fmt.Println("Environment is healthy")
	} else {
		fmt.Println("Environment is unhealthy")
		os.Exit(1)
	}
}

func loadConfig() (*config.Medik, error) {
	content, err := os.ReadFile(ConfigFile)

	if err != nil {
		return nil, err
	}

	return config.Parse(string(content))
}

func loadEnv() (bool, error) {
	if !NoUseEnv {
		return true, nil
	}

	content, err := os.ReadFile(EnvFile)

	if err != nil {
		return false, err
	}

	return useEnv(string(content))
}

func useEnv(content string) (bool, error) {
	env, err := parse.ParseEnvFile(content)

	if err != nil {
		return false, err
	}

	for k, v := range env {
		err := os.Setenv(k, v)

		if err != nil {
			return false, err
		}
	}

	return true, nil
}
