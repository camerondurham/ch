package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/camerondurham/ch/cmd/dockerutil"
	"github.com/camerondurham/ch/cmd/util"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create ENVIRONMENT_NAME {--file DOCKERFILE|--image DOCKER_IMAGE} [--volume PATH_TO_DIRECTORY] [--shell SHELL_CMD]",
	Short: "Create docker environment config",
	Long: `Create docker environment config with new name. 
	Will look for your Dockerfile in the current directory 
	if you do not explicitly set --file.`,
	Args: cobra.MinimumNArgs(1),
	Run:  CreateCmd,
}

// CreateCmd creates a new Docker environment
func CreateCmd(cmd *cobra.Command, args []string) {

	/*
		-1. parse new environment name
			-1.1 get list of current environments
			-1.2 check if "name" exists already
		0. parse file, shell, volume
		1. save new environment in config file
	*/

	name := args[0]
	queryName := viper.GetStringMapString(name)
	if len(queryName) > 0 {
		log.Fatalf("environment name [%s] already exists", name)
	}

	opts, err := parseContainerOpts(cmd, name)

	if err != nil {
		log.Fatal("failed to parse args: ", err)
	}

	cli, err := client.NewEnvClient()
	ctx := context.Background()

	if err != nil {
		log.Fatal("error creating Docker client: are you sure Docker is running?")
	}

	if opts.BuildOpts != nil {
		err = dockerutil.BuildImageWithContext(ctx,
			cli,
			opts.BuildOpts.DockerfilePath,
			opts.BuildOpts.Context,
			opts.BuildOpts.Tag)
	} else {
		err = dockerutil.PullImage(ctx,
			cli,
			opts.PullOpts.ImageName)
	}

	if err != nil {
		log.Fatal("error creating image: ", err)
	}

	util.DebugPrint(fmt.Sprintf("Saving environment: %v", *opts))

	// save new environment opts into config file
	viper.Set(name, *opts)
	viper.WriteConfig()
}
func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringP("file", "f", "Dockerfile", "path to Dockerfile")

	createCmd.Flags().StringP("image", "i", "", "image name to pull from DockerHub")

	createCmd.Flags().StringP("volume", "v", "", "volume to mount to the working directory")

	createCmd.Flags().String("shell", "/bin/sh", "default shell to use when logging into environment")

	createCmd.Flags().String("context", ".", "context to build Dockerfile")

	// Example to bind any other flags to all viper flags
	// viper.BindPFlag("context", createCmd.PersistentFlags().Lookup("context"))
}

var (
	errorCreateImageFieldsNotPresent = errors.New("file or image must be provided to create container")
	errorBuildImageFieldsNotPresent  = errors.New("file and context must be provided to build a container")
)

func parseContainerOpts(cmd *cobra.Command, environmentName string) (*ContainerOpts, error) {
	if file, _ := cmd.Flags().GetString("file"); file != "" {
		if contextDirName, _ := cmd.Flags().GetString("context"); contextDirName != "" {
			volumeName, shellCmd := parseOptional(cmd)
			return &ContainerOpts{
				BuildOpts: &BuildOpts{
					DockerfilePath: file,
					Context:        contextDirName,
					Tag:            environmentName,
				},
				Volume: volumeName,
				Shell:  shellCmd,
			}, nil
		} else {
			return &ContainerOpts{}, errorBuildImageFieldsNotPresent
		}
	}

	if imageName, _ := cmd.Flags().GetString("image"); imageName != "" {
		volumeName, shellCmd := parseOptional(cmd)
		return &ContainerOpts{
			PullOpts: &PullOpts{
				ImageName: imageName,
			},
			Volume: volumeName,
			Shell:  shellCmd,
		}, nil
	}

	return nil, errorCreateImageFieldsNotPresent
}

func parseOptional(cmd *cobra.Command) (volumeName string, shellCmd string) {
	volumeName, _ = cmd.Flags().GetString("volume")
	shellCmd, _ = cmd.Flags().GetString("shell")
	return
}
