package cmd

import (
	"fmt"
	"github.com/camerondurham/ch/cmd/util"
	"github.com/camerondurham/ch/version"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

const (
	RepositoryUrl = "https://github.com/camerondurham/ch"
)

// TODO: cleanup skeleton code
// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "ch",
	Short:   "A simple container helper to create and manage Docker environments",
	Version: version.PkgVersion,
	// 	Long: `A longer description that spans multiple lines and likely contains
	// examples and usage of using your application. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// TODO: only check if user didn't specify they don't want upgrade reminders
	latestVersion, err := util.GetLatestVersion(util.GetRequest)
	if err != nil {
		util.DebugPrint(fmt.Sprintf("ignoring version check since error occured when retrieving latest version: %v\n", err))
	} else if version.PkgVersion != "" && latestVersion != version.PkgVersion {
		fmt.Printf("A new version of ch is available!\n"+
			"You are running version %s but the latest version is %s."+
			"\nSee %s instructions on upgrading.\n",
			version.PkgVersion,
			latestVersion,
			RepositoryUrl)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ch.yaml)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".ch" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".ch")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		util.DebugPrint(fmt.Sprint("Using config file:", viper.ConfigFileUsed()))
	}
}
