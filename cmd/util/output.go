package util

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	printConfigNest0 = "%s:\t%s\n"
	printConfigNest1 = "\t%s:\t%s\n"
)

var (
	ErrDoesNotExist = errors.New("does not exist")
)

func PrintConfig(envName string, opts *ContainerOpts) {
	fmt.Println(buildString(envName, opts))
}

func PrintStartHelpMessage(envName string) {
	fmt.Printf("\nStart this environment with: \n\tch start %s\n", envName)
}

func buildString(envName string, opts *ContainerOpts) string {
	var b strings.Builder

	fmt.Fprintf(&b, printConfigNest0, "Name", envName)
	bo := opts.BuildOpts
	po := opts.PullOpts

	if bo != nil {
		fmt.Fprintf(&b, printConfigNest1, "Dockerfile", bo.DockerfilePath)
		fmt.Fprintf(&b, printConfigNest1, "Context", bo.Context)
		fmt.Fprintf(&b, printConfigNest1, "    Tag", bo.Tag)
	} else if po != nil {
		fmt.Fprintf(&b, printConfigNest1, "Image", po.ImageName)
	}

	if opts.HostConfig != nil {
		if len(opts.HostConfig.Binds) > 0 {
			for _, v := range opts.HostConfig.Binds {
				fmt.Fprintf(&b, printConfigNest1, "Volume", fmt.Sprintf("%s", v))
			}
		}

		if len(opts.HostConfig.SecurityOpt) > 0 {
			for _, v := range opts.HostConfig.SecurityOpt {
				fmt.Fprintf(&b, printConfigNest1, "SecOpt", fmt.Sprintf("%s", v))
			}
		}

		if len(opts.HostConfig.CapAdd) > 0 {
			for _, v := range opts.HostConfig.CapAdd {
				fmt.Fprintf(&b, printConfigNest1, "CapAdd", fmt.Sprintf("%s", v))
			}
		}

		if len(opts.HostConfig.PortBindings) > 0 {
			for k, v := range opts.HostConfig.PortBindings {
				fmt.Fprintf(&b, printConfigNest1, "Port", fmt.Sprintf("%s:%s", v[0].HostPort, k.Port()))
			}
		}
	}

	return b.String()
}

// DebugPrint if DEBUG environment variable is set
func DebugPrint(msg string) {
	if _, ok := os.LookupEnv("DEBUG"); ok {
		fmt.Println(msg)
	}
}
