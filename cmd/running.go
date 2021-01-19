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
	"fmt"
	"github.com/camerondurham/ch/cmd/util"
	"github.com/docker/docker/api/types/filters"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

// runningCmd represents the running command
var runningCmd = &cobra.Command{
	Use:   "running",
	Short: "List running environments",
	Long: `List all running Docker containers created by the container-helper.
Any Docker containers not managed by the container-helper will be ignored.
To see all running containers, run: docker ps`,
	Run: RunningCmd,
}

func RunningCmd(cmd *cobra.Command, args []string) {
	cli, err := util.NewCliClient()
	if err != nil {
		fmt.Printf("error: cannot create new CLI ApiClient: %v\n", err)
		os.Exit(1)
	}

	envs := cli.Containers()
	environmentNames := make([]filters.KeyValuePair, 0)
	for name, _ := range envs {
		environmentNames = append(environmentNames, filters.Arg("name", name))
	}
	f := filters.NewArgs(environmentNames...)
	list, err := cli.DockerClient().GetRunning(f, false)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintf(w, "ENVIRONMENT\tIMAGE NAME\tCREATED\n")
	for _, info := range list {
		now := time.Now().Unix()
		elapsed := now - info.Created
		var timeSinceCreated string
		if elapsed >= 60 {
			timeSinceCreated = fmt.Sprintf("%v minutes ago", elapsed/60)
		} else {
			timeSinceCreated = fmt.Sprintf("%v seconds ago", elapsed)
		}
		fmt.Fprintf(w, "%s\t%s\t%v\n", strings.TrimPrefix(info.Names[0], "/"), info.Image, timeSinceCreated)
	}
	w.Flush()
}

func init() {
	rootCmd.AddCommand(runningCmd)
}
