package cmd

import (
	"context"
	"fmt"
	"github.com/camerondurham/ch/cmd/util"
	"github.com/spf13/cobra"
	"os"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:     "stop ENVIRONMENT_NAME",
	Short:   "Stop a running environment (a running Docker container)",
	Args:    cobra.ExactArgs(1),
	Version: rootCmd.Version,
	Run:     StopCmd,
}

func StopCmd(cmd *cobra.Command, args []string) {
	envName := args[0]
	cli, err := util.NewCliClient()
	if err != nil {
		fmt.Printf("error: cannot create new CLI ApiClient: %v\n", err)
		os.Exit(1)
	}

	envs := cli.Containers()

	if _, ok := envs[envName]; ok {

		c, _ := cli.Container(envName)
		if c == nil {
			fmt.Printf("%v is not running\n", envName)
			os.Exit(1)
		} else {
			ctx := context.Background()

			containerID := c.ID

			err := cli.DockerClient().StopContainer(ctx, containerID, nil)
			if err != nil {
				fmt.Printf("container not running\n")
			} else {
				fmt.Printf("stopped container: %v\n", envName)
				err = cli.DockerClient().RemoveContainer(ctx, envName)
				if err != nil {
					util.DebugPrint(fmt.Sprintf("error removing container: %v", err))
				}
			}

		}

	} else {
		util.PrintEnvNotFoundMsg(envName)
		os.Exit(1)
	}
}
func init() {
	rootCmd.AddCommand(stopCmd)
}
