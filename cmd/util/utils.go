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
	if len(opts.Volume) > 0 {
		for _, v := range opts.Volume {
			fmt.Printf(printConfigNest1, "Volume", fmt.Sprintf(" %s ", v))
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
