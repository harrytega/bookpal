package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"test-project/cmd/db"
	"test-project/cmd/env"
	"test-project/cmd/probe"
	"test-project/cmd/server"
	"test-project/internal/config"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: config.GetFormattedBuildArgs(),
	Use:     "app",
	Short:   config.ModuleName,
	Long: fmt.Sprintf(`%v

A stateless RESTful JSON service written in Go.
Requires configuration through ENV.`, config.ModuleName),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate(`{{printf "%s\n" .Version}}`)

	// attach the subcommands
	rootCmd.AddCommand(
		db.New(),
		env.New(),
		probe.New(),
		server.New(),
	)
}
