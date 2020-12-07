package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/camerondurham/ch/cmd/config"
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
// TODO: restructure for easier testing
// https://stackoverflow.com/questions/35827147/cobra-viper-golang-how-to-test-subcommands
func CreateCmd(cmd *cobra.Command, args []string) {

	name := args[0]
	queryName := viper.GetStringMapString(fmt.Sprintf("envs.%s", name))
	replace, _ := cmd.Flags().GetBool("replace")

	if len(queryName) > 0 && !replace {
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
		err = util.BuildImageWithContext(ctx,
			cli,
			opts.BuildOpts.DockerfilePath,
			opts.BuildOpts.Context,
			opts.BuildOpts.Tag)
	} else {
		err = util.PullImage(ctx,
			cli,
			opts.PullOpts.ImageName)
	}

	if err != nil {
		log.Fatal("cannot create new environment, error creating image: ", err)
	}

	util.DebugPrint(fmt.Sprintf("saving environment: %v", opts))

	envs, err := util.GetEnvs()

	if err != nil {
		if err == util.ErrDoesNotExist {
			// save new environment opts into config file
			viper.Set(fmt.Sprintf("envs.%s", name), opts)
		} else {
			log.Fatal("cannot read environment")
		}
	} else {
		envs[name] = opts
		viper.Set("envs", envs)
	}

	err = viper.WriteConfig()
	if err != nil {
		log.Fatal("error saving config: ", err)
	}

	util.PrintConfig(name, opts)
}
func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringP("file", "f", "Dockerfile", "path to Dockerfile")
	createCmd.Flags().StringP("image", "i", "", "image name to pull from DockerHub")
	createCmd.Flags().StringP("volume", "v", "", "volume to mount to the working directory")
	createCmd.Flags().String("shell", "/bin/sh", "default shell to use when logging into environment")
	createCmd.Flags().String("context", ".", "context to build Dockerfile")
	createCmd.Flags().Bool("replace", false, "replace environment if it already exists")
}

var (
	errorCreateImageFieldsNotPresent = errors.New("file or image must be provided to create container")
	errorBuildImageFieldsNotPresent  = errors.New("file and context must be provided to build a container")
)

func parseContainerOpts(cmd *cobra.Command, environmentName string) (*config.ContainerOpts, error) {
	if file, _ := cmd.Flags().GetString("file"); file != "" {
		if contextDirName, _ := cmd.Flags().GetString("context"); contextDirName != "" {
			volumeName, shellCmd := parseOptional(cmd)
			return &config.ContainerOpts{
				BuildOpts: &config.BuildOpts{
					DockerfilePath: file,
					Context:        contextDirName,
					Tag:            environmentName,
				},
				Volume: volumeName,
				Shell:  shellCmd,
			}, nil
		} else {
			return &config.ContainerOpts{}, errorBuildImageFieldsNotPresent
		}
	}

	if imageName, _ := cmd.Flags().GetString("image"); imageName != "" {
		volumeName, shellCmd := parseOptional(cmd)
		return &config.ContainerOpts{
			PullOpts: &config.PullOpts{
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
