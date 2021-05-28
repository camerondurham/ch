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

var rootCmd = &cobra.Command{
	Use:     "ch",
	Short:   "A simple container helper to create and manage Docker environments",
	Version: version.PkgVersion,
}

func Execute() {

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

func initConfig() {
	if cfgFile != "" {
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

	if err := viper.ReadInConfig(); err == nil {
		util.DebugPrint(fmt.Sprint("Using config file: ", viper.ConfigFileUsed()))
	}
}
