/*
Copyright © 2020 Cameron Durham <cameron.r.durham@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

*/
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
	Short:   "Start a shell in an environment",
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
			util.DebugPrint("starting non-containerIsRunning container because autostart flag used\n")
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
		fmt.Printf("no such environment: %v\n", envName)
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
