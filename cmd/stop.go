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
		fmt.Printf("environment does not exist: %v\n", envName)
		os.Exit(1)
	}
}
func init() {
	rootCmd.AddCommand(stopCmd)
}
