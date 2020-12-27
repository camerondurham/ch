package util

import (
	"errors"
	"fmt"
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
	if len(opts.HostConfig.Binds) > 0 {
		for _, v := range opts.HostConfig.Binds {
			fmt.Printf(printConfigNest1, "Volume", fmt.Sprintf(" %s ", v))
		}
	}

	if len(opts.HostConfig.SecurityOpt) > 0 {
		for _, v := range opts.HostConfig.SecurityOpt {
			fmt.Printf(printConfigNest1, "SecOpt", fmt.Sprintf(" %s ", v))
		}
	}

	if len(opts.HostConfig.CapAdd) > 0 {
		for _, v := range opts.HostConfig.CapAdd {
			fmt.Printf(printConfigNest1, "CapAdd", fmt.Sprintf(" %s ", v))
		}
	}

	fmt.Println()
}

// DebugPrint if DEBUG environment variable is set
func DebugPrint(msg string) {
	if _, ok := os.LookupEnv("DEBUG"); ok {
		fmt.Println(msg)
	}
}
