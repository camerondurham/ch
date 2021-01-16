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
	"os"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:     "start ENVIRONMENT_NAME",
	Short:   "start environment in the background",
	Args:    cobra.ExactArgs(1),
	Version: rootCmd.Version,
	Run:     StartCmd,
}

func StartCmd(cmd *cobra.Command, args []string) {

	envName := args[0]
	if envName == "" {
		fmt.Printf("you must provide an environment name\n")
		_ = cmd.Usage()
		os.Exit(1)
	}

	cli, err := util.NewCliClient()
	if err != nil {
		fmt.Printf("error: cannot create new CLI ApiClient: %v\n", err)
		os.Exit(1)
	}

	envs := cli.Containers()

	if containerOpts, ok := envs[envName]; ok {
		running, err := cli.Running()
		_, ok := running[envName]
		if err == nil && ok {
			fmt.Printf("%v is already running\n", envName)
			os.Exit(0)
		}
		startEnvironment(cli, containerOpts, envName)
	} else {
		fmt.Printf("no such environment: %v\n", envName)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func startEnvironment(client *util.Cli, containerOpts *util.ContainerOpts, envName string) {

	ctx := context.Background()

	imageName := getImageName(envName, containerOpts)
	containerConfig := &container.Config{Image: imageName, Tty: true, AttachStdin: true}
	if containerOpts.Shell != "" {
		containerConfig.Shell = []string{containerOpts.Shell}
	}

	hostConfig := createHostConfig(containerOpts)

	resp, err := client.DockerClient().CreateContainer(ctx,
		containerConfig,
		envName,
		hostConfig)

	if err != nil {
		fmt.Printf("error creating container: %v\n", err)
		os.Exit(1)
	}

	client.DockerClient().StartContainer(ctx, resp.ID)
	fmt.Printf("[%v] started...\n", envName)
	util.DebugPrint(fmt.Sprintf("containerID:\n%v", resp.ID))

	running, _ := client.Running()

	running = updateRunning(running, resp.ID, envName)

	viper.Set("running", running)

	err = viper.WriteConfig()
	if err != nil {
		fmt.Printf("failed saving running containers\n")
	}
}

func updateRunning(running map[string]string, containerID string, envName string) map[string]string {
	if running == nil {
		running = make(map[string]string)
	}
	running[envName] = containerID
	return running
}

func getImageName(envName string, containerOpts *util.ContainerOpts) string {
	if containerOpts.BuildOpts != nil {
		return envName
	} else {
		return containerOpts.PullOpts.ImageName
	}
}

func createHostConfig(containerOpts *util.ContainerOpts) *container.HostConfig {
	if containerOpts.HostConfig != nil {
		return &container.HostConfig{
			Binds:       containerOpts.HostConfig.Binds,
			CapAdd:      containerOpts.HostConfig.CapAdd,
			Privileged:  containerOpts.HostConfig.Privileged,
			SecurityOpt: containerOpts.HostConfig.SecurityOpt,
		}
	} else {
		return &container.HostConfig{}
	}

}
