package cmd

import (
	"fmt"
	"os"

	"github.com/OJarrisonn/medik/pkg/config"
	"github.com/OJarrisonn/medik/pkg/medik"
	"github.com/OJarrisonn/medik/pkg/runner"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     medik.Name,
	Version: medik.Version,
	Short:   medik.Description,
	Args:    cobra.ArbitraryArgs,
	RunE:    run,
}

var ConfigFile string
var EnvFile string
var UseEnv bool

func init() {
	rootCmd.PersistentFlags().StringVarP(&ConfigFile, "config", "c", medik.DefaultConfigFile, "Config file to use")
	rootCmd.PersistentFlags().StringVarP(&EnvFile, "env", "e", medik.DefaultEnvFile, "Env file to use")
	rootCmd.PersistentFlags().BoolVar(&UseEnv, "no-env", medik.DefaultUseEnv, "Won't use an env file")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()

	if err != nil {
		return err
	}

	success, errs := runner.Run(cfg, args)

	for _, e := range errs {
		fmt.Println(e)
	}

	if success {
		fmt.Println("Environment is healthy")
		return nil
	} else {
		fmt.Println("Environment is unhealthy")
		os.Exit(1)
		return nil
	}
}

func loadConfig() (*config.Medik, error) {
	content, err := os.ReadFile(ConfigFile)

	if err != nil {
		return nil, err
	}

	return config.Parse(string(content))
}
