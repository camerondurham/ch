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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

// deleteCmd represents the delete command
var (
	deleteCmd = &cobra.Command{
		Use:     "delete ENVIRONMENT_NAME",
		Short:   "deletes a given config",
		Args:    cobra.MinimumNArgs(1),
		Version: rootCmd.Version,
		Run:     DeleteCmd,
	}
)

func DeleteCmd(cmd *cobra.Command, args []string) {
	envName := args[0]
	envs, err := util.GetEnvs()
	if err != nil {
		log.Fatalf("cannot read config: %v", err)
	}

	cli, err := util.NewCliClient()
	if err != nil {
		fmt.Printf("error: cannot create new CLI ApiClient: %v\n", err)
		os.Exit(1)
	}

	if envs, err := removeEnvironment(envName, envs, cli); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	} else {
		viper.Set("envs", envs)
		err = viper.WriteConfig()
		if err != nil {
			fmt.Printf("cannot write to config: %v\n", err)
			os.Exit(1)
		} else {
			fmt.Printf("environment deleted: %v\n", envName)
		}
	}
}
func init() {
	rootCmd.AddCommand(deleteCmd)
}

func removeEnvironment(envName string, envs map[string]*util.ContainerOpts, cli util.ContainerClient) (map[string]*util.ContainerOpts, error) {
	if _, ok := envs[envName]; ok {

		if cli.ContainerIsRunning(envName) {
			return envs, fmt.Errorf("environment [%s] is currently running! \nPlease stop the environment before deleting it.\n", envName)
		}

		delete(envs, envName)
		return envs, nil
	} else {
		return envs, fmt.Errorf("environment [%s] not found", envName)
	}
}
