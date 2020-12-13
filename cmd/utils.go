package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

const (
	printConfigNest0 = "%s:\t%s\n"
	printConfigNest1 = "\t%s:\t%s\n"
)

var (
	ErrDoesNotExist = errors.New("does not exist")
)

func PrintConfig(envName string, opts *ContainerOpts) {
	fmt.Printf(printConfigNest0, "Name", envName)
	bo := opts.BuildOpts
	po := opts.PullOpts
	if bo != nil {
		fmt.Printf(printConfigNest1, "Dockerfile", bo.DockerfilePath)
		fmt.Printf(printConfigNest1, "Context", bo.Context)
		fmt.Printf(printConfigNest1, "    Tag", bo.Tag)
	} else {
		fmt.Printf(printConfigNest1, "Image", po.ImageName)
	}
	fmt.Println()
}

// DebugPrint if DEBUG environment variable is set
func DebugPrint(msg string) {
	if _, ok := os.LookupEnv("DEBUG"); ok {
		fmt.Println(msg)
	}
}

func GetEnvs() (envs map[string]*ContainerOpts, err error) {
	if !viper.IsSet("envs") {
		return nil, ErrDoesNotExist
	}

	envs = make(map[string]*ContainerOpts)
	err = viper.UnmarshalKey("envs", &envs)
	return
}

func GetEnvsOrDie() map[string]*ContainerOpts {
	envs, err := GetEnvs()
	if err != nil {
		fmt.Printf("error retrieving envs: %v\n", err)
		os.Exit(1)
	}
	return envs
}

func GetRunning() (running map[string]string, err error) {
	if !viper.IsSet("running") {
		return nil, ErrDoesNotExist
	}

	running = make(map[string]string)
	err = viper.UnmarshalKey("running", &running)
	return
}
