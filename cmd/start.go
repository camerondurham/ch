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
	"github.com/docker/docker/api/types/container"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start environment in the background",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		envName := args[0]

		cli, err := util.NewCliClient()
		if err != nil {
			fmt.Printf("error: cannot create new CLI ApiClient: %v", err)
			os.Exit(1)
		}

		envs := cli.Containers()

		if containerOpts, ok := envs[envName]; ok {
			StartEnvironment(cli, containerOpts, envName)
		} else {
			fmt.Printf("no such environment: %v", envName)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func StartEnvironment(client *util.Cli, containerOpts *util.ContainerOpts, envName string) {

	ctx := context.Background()

	containerConfig := &container.Config{Image: envName, Tty: true, AttachStdin: true}
	if containerOpts.Shell != "" {
		containerConfig.Shell = []string{containerOpts.Shell}
	}

	if len(containerOpts.Volume) > 0 {
		// TODO(cadurham): error check attaching volumes
		containerConfig.Volumes = containerOpts.Volume
	}

	resp := util.CreateContainer(ctx, client.Client(), containerConfig, envName)

	util.StartContainer(ctx, client.Client(), resp.ID)
	fmt.Printf("started... imageID: %v", resp.ID)

	running, _ := client.Running()

	if running == nil {
		running = make(map[string]string)
	}

	running[envName] = resp.ID
	viper.Set("running", running)

	err := viper.WriteConfig()
	if err != nil {
		log.Printf("failed saving running containers")
	}

}
