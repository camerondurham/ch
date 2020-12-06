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
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
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
		if !viper.IsSet("envs") {
			fmt.Print("no environments exist")
			os.Exit(1)
		}
		envs, err := util.GetEnvs()
		if err != nil {
			fmt.Printf("error retrieving envs: %v\n", err)
			os.Exit(1)
		}

		if containerOpts, ok := envs[envName]; ok {
			ctx := context.Background()
			cli, err := client.NewEnvClient()

			containerConfig := &container.Config{Image: envName}
			if containerOpts.Shell != "" {
				containerConfig.Shell = []string{containerOpts.Shell}
			}
			if err != nil {
				log.Fatal("error creating Docker client: are you sure Docker is running?")
			}
			resp := util.CreateContainer(ctx, cli, &container.Config{
				Image: envName,
			}, envName)

			util.StartContainer(ctx, cli, resp.ID)

			fmt.Printf("started image: %v", resp.ID)
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
