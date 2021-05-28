package cmd

import (
	"context"
	"fmt"
	"github.com/camerondurham/ch/cmd/util"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/errdefs"
	"github.com/spf13/cobra"
	"os"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:     "start ENVIRONMENT_NAME",
	Short:   "start a Docker environment in the background",
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
		if cli.ContainerIsRunning(envName) {
			fmt.Printf("%v is already running\n", envName)
			os.Exit(0)
		}
		startEnvironment(cli, containerOpts, envName)
	} else {
		util.PrintEnvNotFoundMsg(envName)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func startEnvironment(client util.ContainerClient, containerOpts *util.ContainerOpts, envName string) (containerID string) {

	ctx := context.Background()

	imageName := getImageName(envName, containerOpts)
	containerConfig := &container.Config{Image: imageName, Tty: true, AttachStdin: true}
	if containerOpts.Shell != "" {
		containerConfig.Shell = []string{containerOpts.Shell}
	}

	hostConfig := createHostConfig(containerOpts)

	retry := 3
	success := false
	var resp container.ContainerCreateCreatedBody
	var err error

	for i := 0; !success && i < retry; i++ {
		resp, err = client.DockerClient().CreateContainer(ctx,
			containerConfig,
			envName,
			hostConfig)

		if err != nil {

			if errdefs.IsConflict(err) {
				// handle conflicting containers by removing them
				err = client.DockerClient().RemoveContainer(ctx, envName)
				if err != nil {
					util.DebugPrint(fmt.Sprintf("error removing container: %v\n", err))
				}
			} else if errdefs.IsNotFound(err) {
				// handle if user has removed image by rebuilding or re-pulling image
				err = initializeImage(ctx, client, containerOpts)
				if err != nil {
					util.DebugPrint(fmt.Sprintf("error recreating or pulling image: %v\n", err))
				}
			} else {
				fmt.Printf("unexpected error creating container: %v\n", err)
				os.Exit(1)
			}

		} else {
			success = true
		}
	}

	err = client.DockerClient().StartContainer(ctx, resp.ID)
	if err != nil {
		fmt.Printf("WARNING: error encountered starting container\n")
		util.DebugPrint(fmt.Sprintf("%v", err))
	}
	fmt.Printf("[%v] started...\n", envName)
	util.DebugPrint(fmt.Sprintf("containerID:\n%v\n", resp.ID))

	containerID = resp.ID
	return
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
			Binds:        containerOpts.HostConfig.Binds,
			CapAdd:       containerOpts.HostConfig.CapAdd,
			Privileged:   containerOpts.HostConfig.Privileged,
			PortBindings: containerOpts.HostConfig.PortBindings,
			SecurityOpt:  containerOpts.HostConfig.SecurityOpt,
		}
	} else {
		return &container.HostConfig{}
	}

}
