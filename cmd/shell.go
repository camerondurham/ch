package cmd

import (
	"context"
	"fmt"
	"github.com/camerondurham/ch/cmd/util"
	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
	"os"
)

const (
	autostartFlag      = "force-start"
	autostartFlagShort = "f"
)

// shellCmd represents the shell command
var shellCmd = &cobra.Command{
	Use:     "shell ENVIRONMENT_NAME",
	Short:   "get a shell in a running Docker environment",
	Args:    cobra.ExactArgs(1),
	Version: rootCmd.Version,
	Run:     ShellCmd,
}

func ShellCmd(cmd *cobra.Command, args []string) {
	envName := args[0]
	autostart, _ := cmd.Flags().GetBool(autostartFlag)

	cli, err := util.NewCliClient()
	if err != nil {
		fmt.Printf("error: cannot create new CLI ApiClient: %v\n", err)
		os.Exit(1)
	}

	envs := cli.Containers()

	if containerOpts, ok := envs[envName]; ok {

		var containerID string

		containerIsRunning := cli.ContainerIsRunning(envName)

		if !autostart && !containerIsRunning {
			fmt.Printf(getNotRunningMsg(envName))
			os.Exit(1)
		} else if autostart && !containerIsRunning {
			util.DebugPrint("starting container because autostart flag used\n")
			containerID = startEnvironment(cli, containerOpts, envName)
		} else {
			c, err := cli.Container(envName)
			if err != nil {
				fmt.Printf("error getting container for environment %s: %v", envName, err)
				return
			}
			containerID = c.ID
		}

		util.DebugPrint(fmt.Sprintf("starting container: %v\n", containerID))

		err = cli.ApiClient().CreateExecInteractive(context.Background(), cli, containerID, types.ExecConfig{
			Cmd:          []string{containerOpts.Shell},
			Tty:          true,
			AttachStdin:  true,
			AttachStderr: true,
			AttachStdout: true,
		})

		if err != nil {
			fmt.Printf("error creating shell: %v\n", err)
		}

	} else {
		util.PrintEnvNotFoundMsg(envName)
		os.Exit(1)
	}
}
func init() {
	rootCmd.AddCommand(shellCmd)
	shellCmd.Flags().BoolP(autostartFlag, autostartFlagShort, false, "autostart the environment if not running")
}

func getNotRunningMsg(envName string) string {
	return fmt.Sprintf(`%v is not running, please run: 
	ch start %v

or start container automatically with:
	ch shell %v --%v

`, envName, envName, envName, autostartFlag)
}
