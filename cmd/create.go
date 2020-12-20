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
	"github.com/camerondurham/ch/cmd/util"
	"log"
	"strings"

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

	cli, err := util.NewCliClient()
	ctx := context.Background()

	if err != nil {
		log.Printf("error: %v", err)
		log.Fatal("error creating Docker client: are you sure Docker is running?")
	}

	if opts.BuildOpts != nil {
		err = util.BuildImageWithContext(ctx,
			cli.Client(),
			opts.BuildOpts.DockerfilePath,
			opts.BuildOpts.Context,
			opts.BuildOpts.Tag)
	} else {
		err = util.PullImage(ctx,
			cli.Client(),
			opts.PullOpts.ImageName)
	}

	if err != nil {
		log.Fatal("cannot create new environment, error creating image: ", err)
	}

	util.DebugPrint(fmt.Sprintf("saving environment: %v", opts))

	err = viper.ReadInConfig()
	if err != nil {
		util.DebugPrint(fmt.Sprintf("error reading in config: %v", err))
	}
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

	util.SetEnvs(envs)
	err = viper.WriteConfig()

	if err != nil {
		log.Print("error saving config: ", err)
	}

	util.PrintConfig(name, opts)
}
func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringP("file", "f", "Dockerfile", "path to Dockerfile")
	createCmd.Flags().StringP("image", "i", "", "image name to pull from DockerHub")
	createCmd.Flags().StringArrayP("volume", "v", nil, "volume to mount to the working directory")
	createCmd.Flags().String("shell", "/bin/sh", "default shell to use when logging into environment")
	createCmd.Flags().String("context", ".", "context to build Dockerfile")
	createCmd.Flags().Bool("replace", false, "replace environment if it already exists")
}

var (
	errorCreateImageFieldsNotPresent = errors.New("file or image must be provided to create container")
	errorBuildImageFieldsNotPresent  = errors.New("file and context must be provided to build a container")
)

func parseContainerOpts(cmd *cobra.Command, environmentName string) (*util.ContainerOpts, error) {
	if file, _ := cmd.Flags().GetString("file"); file != "" {
		if contextDirName, _ := cmd.Flags().GetString("context"); contextDirName != "" {
			volumeName, shellCmd := parseOptional(cmd)
			return &util.ContainerOpts{
				BuildOpts: &util.BuildOpts{
					DockerfilePath: file,
					Context:        contextDirName,
					Tag:            environmentName,
				},
				Volume: volumeName,
				Shell:  shellCmd,
			}, nil
		} else {
			return &util.ContainerOpts{}, errorBuildImageFieldsNotPresent
		}
	}

	if imageName, _ := cmd.Flags().GetString("image"); imageName != "" {
		volumeName, shellCmd := parseOptional(cmd)
		return &util.ContainerOpts{
			PullOpts: &util.PullOpts{
				ImageName: imageName,
			},
			Volume: volumeName,
			Shell:  shellCmd,
		}, nil
	}

	return nil, errorCreateImageFieldsNotPresent
}

func parseOptional(cmd *cobra.Command) (volumeName map[string]struct{}, shellCmd string) {
	volNames, _ := cmd.Flags().GetStringArray("volume")
	if len(volNames) > 0 {
		volumeName = make(map[string]struct{})
		for i := 0; i < len(volNames); i++ {
			// TODO: use more sophisticated, less expensive parsing
			arr := strings.Split(volNames[i], ":")
			if len(arr) != 2 {
				fmt.Printf("ignoring invalid volume syntax: %v", volNames[i])
				continue
			}
			volumeName[volNames[i]] = struct{}{}
		}
	}
	shellCmd, _ = cmd.Flags().GetString("shell")
	return
}
