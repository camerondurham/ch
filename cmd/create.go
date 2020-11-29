package cmd

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/camerondurham/ch/cmd/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	debug = true
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create ENVIRONMENT_NAME [--file DOCKERFILE] [--volume PATH_TO_DIRECTORY] [--shell SHELL_CMD]",
	Short: "Create docker environment config",
	Long: `Create docker environment config with new name. 
	Will look for your Dockerfile in the current directory 
	if you do not explicitly set --file.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		/*
			1. parse new environment name
				1.1 get list of current environments
				1.2 check if "name" exists already
			2. parse file, shell, volume
			3. save new environment in config file
		*/

		// TODO: implement with dockerutils

		name := args[0]
		queryName := viper.GetStringMapString(name)
		if len(queryName) > 0 {
			log.Fatalf("environment name [%s] already exists", name)
		}

		file, err := cmd.Flags().GetString("file")
		volume, err := cmd.Flags().GetString("volume")
		shell, err := cmd.Flags().GetString("shell")
		context, err := cmd.Flags().GetString("context")
		tag := name

		command := []string{"docker", "build", "--file", file, "-t", tag, context}
		dockerCmd := exec.Command(command[0], command[1:]...)

		output, err := dockerCmd.CombinedOutput()

		if err != nil {
			log.Fatalf("error creating docker environment: %v", err)
		}

		util.DebugPrint(fmt.Sprintf("output: %s", output))

		opts := ContainerOpts{
			BuildInfo: BuildOpts{
				DockerfilePath: file,
				Context:        context,
				Tag:            tag,
			},
			Volume: volume,
			Shell:  shell,
		}

		if debug {
			log.Print("created docker environment")
			log.Printf("Saving environment: %s", opts)
		}

		viper.Set(name, opts)

		viper.WriteConfig()

	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.PersistentFlags().StringP("file", "f", "Dockerfile", "path to Dockerfile")

	createCmd.PersistentFlags().StringP("volume", "v", "", "volume to mount to the working directory")

	createCmd.PersistentFlags().String("shell", "/bin/sh", "default shell to use when logging into environment")

	createCmd.PersistentFlags().String("context", ".", "context to build Dockerfile")

	// Example to bind any other flags to all viper flags
	// viper.BindPFlag("context", createCmd.PersistentFlags().Lookup("context"))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
