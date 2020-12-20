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
	"errors"
	"fmt"
	"log"

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
		err = BuildImageWithContext(ctx,
			cli,
			opts.BuildOpts.DockerfilePath,
			opts.BuildOpts.Context,
			opts.BuildOpts.Tag)
	} else {
		err = PullImage(ctx,
			cli,
			opts.PullOpts.ImageName)
	}

	if err != nil {
		log.Fatal("cannot create new environment, error creating image: ", err)
	}

	DebugPrint(fmt.Sprintf("saving environment: %v", opts))

	envs, err := GetEnvs()

	if err != nil {
		if err == ErrDoesNotExist {
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

	PrintConfig(name, opts)
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
