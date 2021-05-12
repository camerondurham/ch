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
