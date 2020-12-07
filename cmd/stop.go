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
	"context"
	"fmt"
	"github.com/camerondurham/ch/cmd/util"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop ENVIRONMENT_NAME",
	Short: "Stop a running environment (a running Docker container)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		envName := args[0]
		envs := util.GetEnvsOrDie()

		if _, ok := envs[envName]; ok {

			running, err := util.GetRunning()
			containerID, ok := running[envName]
			if err == util.ErrDoesNotExist || !ok {
				fmt.Printf("%v is not running", envName)
				os.Exit(1)
			} else {

				// TODO: make helper function to do this
				ctx := context.Background()
				cli, err := client.NewEnvClient()
				if err != nil {
					// yikes! 3 layers of nesting!
					log.Fatal("error creating Docker client: are you sure Docker is running?")
				}

				util.StopContainer(ctx, cli, containerID, nil)
				util.RemoveContainer(ctx, cli, envName)

				// TODO: use standard text formatting for all errors, look for library?
				fmt.Printf("stopped container: %v", envName)

				delete(running, envName)

				// TODO: use other storage for running containers, possibly discover through Docker API?
				viper.Set("running", running)
				err = viper.WriteConfig()
				if err != nil {
					fmt.Printf("error writing changes to config")
				}

			}

		} else {
			fmt.Printf("environment does not exist: %v", envName)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
