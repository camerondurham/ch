package cmd

import (
	"fmt"
	"github.com/camerondurham/ch/cmd/util"
	"github.com/spf13/cobra"
	"os"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list [ENVIRONMENT_NAME]",
	Short:   "list configuration for existing environments",
	Version: rootCmd.Version,
	Run:     ListCmd,
}

func ListCmd(cmd *cobra.Command, args []string) {

	util.CheckLatestVersion()

	envName := ""
	if len(args) > 0 {
		envName = args[0]
	}

	configEnvs, err := util.GetEnvs()

	if err != nil {
		if err == util.ErrDoesNotExist {
			fmt.Printf("no environments created yet\n")
			os.Exit(0)
		}
		fmt.Printf("unable to decode config file: %v\nplease check formatting of config and delete if needed\n", err)
		os.Exit(1)
	}

	if envName != "" {
		if v, ok := configEnvs[envName]; ok {
			util.PrintConfig(envName, v)
		} else {
			fmt.Printf("no environment found\n")
		}
	} else {
		for k, v := range configEnvs {
			util.PrintConfig(k, v)
		}
	}

}
func init() {
	rootCmd.AddCommand(listCmd)
}
