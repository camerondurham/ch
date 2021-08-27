package cmd

import (
	"context"
	"fmt"
	"github.com/camerondurham/ch/cmd/util"
	"github.com/spf13/cobra"
	"os"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:     "update [ENVIRONMENT_NAME]",
	Short:   "download or rebuild environment's image",
	Args:    cobra.MinimumNArgs(1),
	Version: rootCmd.Version,
	Run:     UpdateCmd,
}

func UpdateCmd(cmd *cobra.Command, args []string) {
	envName := args[0]
	if envName == "" {
		fmt.Printf("you must provide an environment name\n")
		_ = cmd.Usage()
		os.Exit(1)
	}

	cli, err := util.NewCliClient()

	if err != nil {
		util.PrintDockerClientStartupError(err)
		os.Exit(1)
	}

	envs := cli.Containers()
	if containerOpts, ok := envs[envName]; ok {
		if cli.ContainerIsRunning(envName) {
			fmt.Printf("WARNING: %s is already running\nPlease stop the environment before updating\n", envName)
			os.Exit(1)
		}
		err = util.BuildOrPullContainerImage(context.Background(), cli, containerOpts)
		if err != nil {
			fmt.Printf("error updating environment:\n%v\n", err)
			os.Exit(1)
		}
	} else {
		util.PrintEnvNotFoundMsg(envName)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
