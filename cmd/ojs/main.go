package main

import (
	"fmt"
	"log"

	"github.com/maxuanquang/ojs/internal/configs"
	"github.com/maxuanquang/ojs/internal/utils"
	"github.com/maxuanquang/ojs/internal/wiring"
	"github.com/spf13/cobra"
)

var (
	version    string
	commitHash string
)

const (
	flagConfigFilePath = "config-file-path"
)

func standaloneServer() *cobra.Command {
	command := &cobra.Command{
		Use: "standalone-server",
		RunE: func(cmd *cobra.Command, args []string) error {
			configFilePath, err := cmd.Flags().GetString(flagConfigFilePath)
			if err != nil {
				return err
			}

			app, cleanup, err := wiring.InitializeStandaloneServer(configs.ConfigFilePath(configFilePath), utils.Arguments{})
			if err != nil {
				return err
			}
			defer cleanup()

			app.Start()
			return nil
		},
	}

	command.Flags().String(flagConfigFilePath, "", "If provided, will use the provided config file")

	return command
}

func httpServer() *cobra.Command {
	command := &cobra.Command{
		Use: "http-server",
		RunE: func(cmd *cobra.Command, args []string) error {
			configFilePath, err := cmd.Flags().GetString(flagConfigFilePath)
			if err != nil {
				return err
			}

			app, cleanup, err := wiring.InitializeHTTPServer(configs.ConfigFilePath(configFilePath), utils.Arguments{})
			if err != nil {
				return err
			}
			defer cleanup()

			app.Start()
			return nil
		},
	}

	command.Flags().String(flagConfigFilePath, "", "If provided, will use the provided config file")

	return command
}

func worker() *cobra.Command {
	command := &cobra.Command{
		Use: "worker",
		RunE: func(cmd *cobra.Command, args []string) error {
			configFilePath, err := cmd.Flags().GetString(flagConfigFilePath)
			if err != nil {
				return err
			}

			app, cleanup, err := wiring.InitializeWorker(configs.ConfigFilePath(configFilePath), utils.Arguments{})
			if err != nil {
				return err
			}
			defer cleanup()

			app.Start()
			return nil
		},
	}

	command.Flags().String(flagConfigFilePath, "", "If provided, will use the provided config file")

	return command
}

func cron() *cobra.Command {
	command := &cobra.Command{
		Use: "cron",
		RunE: func(cmd *cobra.Command, args []string) error {
			configFilePath, err := cmd.Flags().GetString(flagConfigFilePath)
			if err != nil {
				return err
			}

			app, cleanup, err := wiring.InitializeCron(configs.ConfigFilePath(configFilePath), utils.Arguments{})
			if err != nil {
				return err
			}
			defer cleanup()

			app.Start()
			return nil
		},
	}

	command.Flags().String(flagConfigFilePath, "", "If provided, will use the provided config file")

	return command
}

func main() {
	rootCommand := &cobra.Command{
		Version: fmt.Sprintf("%s-%s", version, commitHash),
	}
	rootCommand.AddCommand(
		standaloneServer(),
		httpServer(),
		worker(),
		cron(),
	)

	if err := rootCommand.Execute(); err != nil {
		log.Panic(err)
	}
}
