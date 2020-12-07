/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
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
	Use:   "shell",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		envName := args[0]
		autostart, _ := cmd.Flags().GetBool(autostartFlagShort)

		// TODO: create helper function for envName existing
		envs := util.GetEnvsOrDie()
		if containerOpts, ok := envs[envName]; ok {
			ctx, cli := util.DockerClientInitOrDie()
			running, err := util.GetRunning()
			containerID, ok := running[envName]
			if !autostart && (err == util.ErrDoesNotExist || !ok) {
				fmt.Printf(getNotRunningMsg(envName))
				os.Exit(1)
			} else if err == util.ErrDoesNotExist || !ok {
				// TODO: start container
				fmt.Print("error: force starting not implemented yet")
				os.Exit(1)
			}

			util.DebugPrint(fmt.Sprintf("starting container: %v", containerID))
			resp := util.CreateExec(ctx, cli, containerID, types.ExecConfig{
				Detach:       false,
				Tty:          true,
				AttachStdin:  true,
				AttachStdout: true,
				AttachStderr: true,
				Cmd:          []string{containerOpts.Shell},
			}, containerOpts.Shell)

			util.DebugPrint(fmt.Sprintf("started exec: %v\n", resp))

		} else {
			fmt.Printf("no such environment: %v", envName)
		}
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
	shellCmd.Flags().BoolP(autostartFlag, autostartFlagShort, false, "autostart the environment if not running")
}

func getNotRunningMsg(envName string) string {
	return fmt.Sprintf(`%v is not running, please run: 
	ch create %v

or start container automatically with:
	ch shell %v %v\n\n`, envName, envName, envName, autostartFlag)
}
