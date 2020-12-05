package util

import (
	"fmt"
	"github.com/camerondurham/ch/cmd/config"
	"os"
)

const (
	printConfigNest0 = "%s: %s\n"
	printConfigNest1 = "\t%s: %s\n"
)

func PrintConfig(envName string, opts *config.ContainerOpts) {
	fmt.Printf(printConfigNest0, "Environment Name", envName)
	bo := opts.BuildOpts
	po := opts.PullOpts
	if bo != nil {
		fmt.Printf(printConfigNest1, "Dockerfile", bo.DockerfilePath)
		fmt.Printf(printConfigNest1, "Context", bo.Context)
		fmt.Printf(printConfigNest1, "Tag", bo.Tag)
	} else {
		fmt.Printf(printConfigNest1, "Image", po.ImageName)
	}
}

// DebugPrint if DEBUG environment variable is set
func DebugPrint(msg string) {
	if _, ok := os.LookupEnv("DEBUG"); ok {
		fmt.Println(msg)
	}
}
