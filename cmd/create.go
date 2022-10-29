package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/camerondurham/ch/cmd/util"
	"github.com/docker/go-connections/nat"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var (
	createCmd = &cobra.Command{
		Use:   "create ENVIRONMENT_NAME [FLAGS] {--file DOCKERFILE |--image DOCKER_IMAGE } [OPTIONS]",
		Short: "create docker environment config",
		Long: `Create docker environment config with new name.

  You can use the following flag to replace an environment name if it already exists:
    --replace                           Replace any existing environment with the same name

  Will look for your Dockerfile in the current directory
  if you do not explicitly set --file.

  To create environment from a Dockerfile, use:

    --file DOCKERFILE                    Path to Dockerfile. If context is used, filepath must be relative to that path.
    --context PATH                       Context to use when building the Docker image

  To create an environment from a pre-built Docker image, use:

  	--image DOCKER_IMAGE

  You can use the following options:

    --volume HOST_PATH:CONTAINER_PATH    Bind mount a volume (e.g. $PWD:/home )
    --shell PATH                         Command to run for shell (e.g. /bin/sh, /bin/bash)
    --cap-add CAPABILITY                 Add Linux capability (e.g. SYS_PTRACE)
    --security-opt	OPT			             Add security configuration (e.g. "seccomp=unconfined")
    --port NUMBER                        Expose port from container to host
    --privileged                         Give container extended privileges on host (use carefully!)
`,
		Args:    cobra.MinimumNArgs(1),
		Version: rootCmd.Version,
		Run:     CreateCmd,
	}
	errorCreateImageFieldsNotPresent = errors.New("file or image must be provided to create container")
	errorBuildImageFieldsNotPresent  = errors.New("file and context must be provided to build a container")
)

type commandFlags struct {
	file         string
	image        string
	volume       []string
	shell        string
	context      string
	capAdd       []string
	securityOpt  []string
	portBindings []string
	privileged   bool
	replace      bool
}

// CreateCmd creates a new Docker environment
func CreateCmd(cmd *cobra.Command, args []string) {

	name := args[0]
	file, _ := cmd.Flags().GetString("file")
	imageName, _ := cmd.Flags().GetString("image")
	volumeArgs, _ := cmd.Flags().GetStringArray("volume")
	shellCmdArgs, _ := cmd.Flags().GetString("shell")
	contextDir, _ := cmd.Flags().GetString("context")
	capAddArgs, _ := cmd.Flags().GetStringArray("cap-add")
	secOptArgs, _ := cmd.Flags().GetStringArray("security-opt")
	portBindingArgs, _ := cmd.Flags().GetStringArray("port")
	privileged, _ := cmd.Flags().GetBool("privileged")
	replace, _ := cmd.Flags().GetBool("replace")

	cmdFlags := &commandFlags{
		file:         file,
		image:        imageName,
		volume:       volumeArgs,
		shell:        shellCmdArgs,
		context:      contextDir,
		capAdd:       capAddArgs,
		securityOpt:  secOptArgs,
		portBindings: portBindingArgs,
		privileged:   privileged,
		replace:      replace,
	}

	queryName := viper.GetStringMapString(fmt.Sprintf("envs.%s", name))

	if len(queryName) > 0 && !replace {
		log.Fatalf("environment name [%s] already exists", name)
	}

	cli, err := util.NewCliClient()
	ctx := context.Background()

	if err != nil {
		fmt.Printf("error creating Docker client: are you sure Docker is running?\n")
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

	opts, err := parseContainerOpts(name, cli.Validator(), cmdFlags)

	if err != nil {
		fmt.Printf("failed to parse args: %v\n", err)
		os.Exit(1)
	}

	err = util.BuildOrPullContainerImage(ctx, cli, opts)

	if err != nil {
		fmt.Printf("Cannot create new environment. Error from Dockerfile build:\n%v\n", err)
		os.Exit(1)
	}

	util.DebugPrint(fmt.Sprintf("saving environment: %v", opts))
	err = viper.ReadInConfig()
	if err != nil {
		util.DebugPrint(fmt.Sprintf("error reading in config: %v", err))
	}
	envs, err := util.GetEnvs()

	if err != nil {
		if err == util.ErrDoesNotExist {
			// save new environment opts into config file
			envs = make(map[string]*util.ContainerOpts)
		} else {
			util.DebugPrint(fmt.Sprintf("error reading environment: %v", err))
			fmt.Print("cannot read environment from .ch.yaml")
			os.Exit(1)
		}
	}

	envs[name] = opts
	util.SetEnvs(envs)

	if err = viper.SafeWriteConfig(); err != nil {
		// create config file if it doesn't exist
		if os.IsNotExist(err) {
			home, _ := homedir.Dir()
			path := filepath.Join(home, ".ch.yaml")
			err = viper.WriteConfigAs(path)
		}
	}

	err = viper.WriteConfig()

	if err != nil {
		log.Print("error saving config: ", err)
	}

	util.PrintConfig(name, opts)
	util.PrintStartHelpMessage(name)
}
func init() {
	rootCmd.AddCommand(createCmd)

	// Docker build options
	createCmd.Flags().StringP("file", "f", "", "path to Dockerfile, should be relative to context flag")
	createCmd.Flags().String("context", ".", "context to build Dockerfile")

	// Docker pull options
	createCmd.Flags().StringP("image", "i", "", "image name to pull from DockerHub")

	// Docker run options
	createCmd.Flags().StringArrayP("volume", "v", nil, "volume to mount to the working directory")
	createCmd.Flags().StringArrayP("port", "p", nil, "bind host port(s) to container")
	createCmd.Flags().String("shell", "/bin/sh", "default shell to use when logging into environment")
	createCmd.Flags().StringArray("cap-add", nil, "special capacity to add to Docker Container (syscalls)")
	createCmd.Flags().StringArray("security-opt", nil, "security options")
	createCmd.Flags().Bool("privileged", false, "run container as privileged (full root/admin access)")

	// ch options
	createCmd.Flags().Bool("replace", false, "replace environment if it already exists")
}

func parseContainerOpts(environmentName string, v util.Validate, cmdFlags *commandFlags) (*util.ContainerOpts, error) {

	hostConfig, shellCmd := parseHostConfig(cmdFlags.shell, cmdFlags.privileged, cmdFlags.capAdd, cmdFlags.securityOpt, v, cmdFlags.volume, cmdFlags.portBindings)

	if cmdFlags.file != "" {
		if cmdFlags.context != "" {
			contextAbsPath := v.GetAbs(cmdFlags.context)
			return &util.ContainerOpts{
				BuildOpts: &util.BuildOpts{
					DockerfilePath: cmdFlags.file,
					Context:        contextAbsPath,
					Tag:            environmentName,
				},
				HostConfig: hostConfig,
				Shell:      shellCmd,
			}, nil
		} else {
			return nil, errorBuildImageFieldsNotPresent
		}
	} else if cmdFlags.image != "" {
		return &util.ContainerOpts{
			PullOpts: &util.PullOpts{
				ImageName: cmdFlags.image,
			},
			HostConfig: hostConfig,
			Shell:      shellCmd,
		}, nil
	}

	return nil, errorCreateImageFieldsNotPresent
}

func parseHostConfig(shellCmdArg string, privileged bool, capAddArgs []string, secOptArgs []string, v util.Validate, volNameArgs []string, portOpts []string) (hostConfig *util.HostConfig, shellCmd string) {
	hostConfig = &util.HostConfig{}

	if len(volNameArgs) > 0 {
		volumeNames := make([]string, 0)
		for i := 0; i < len(volNameArgs); i++ {
			absPath, err := parseHostContainerPath(volNameArgs[i], v)
			if err != nil {
				fmt.Printf("error parsing mount: %v\ncheck that you have provided a valid path", err)
			} else {
				volumeNames = append(volumeNames, absPath)
			}
		}
		hostConfig.Binds = volumeNames
	}

	if len(portOpts) > 0 {
		ports, portBindings, err := nat.ParsePortSpecs(portOpts)
		if err != nil {
			fmt.Printf("error parsing ports: %v\nerror: %v\n", ports, err)
		}

		hostConfig.PortBindings = portBindings
	}

	if len(capAddArgs) > 0 {
		hostConfig.CapAdd = capAddArgs
	}

	if len(secOptArgs) > 0 {
		hostConfig.SecurityOpt = secOptArgs
	}

	if privileged {
		hostConfig.Privileged = privileged
	}

	if shellCmdArg != "" {
		shellCmd = shellCmdArg
	}
	return
}

func parseHostContainerPath(pathStr string, v util.Validate) (hostContainerAbsPath string, err error) {
	idx := strings.LastIndex(pathStr, ":")

	if idx > 0 {
		hostPath := pathStr[:idx]
		containerPath := pathStr[idx+1:]

		// save original host path to help user debug possible issues
		originalHostPath := hostPath
		hostPath = v.GetAbs(hostPath)

		if idx >= len(pathStr)-1 {
			return "", errors.New("no container path")
		} else if !v.ValidPath(hostPath) {
			return "", errors.New(fmt.Sprintf("invalid host path [%v]", originalHostPath))
		} else {
			hostContainerAbsPath = fmt.Sprintf("%s:%s", hostPath, containerPath)
			return hostContainerAbsPath, nil
		}
	} else {
		return "", errors.New("no container path")
	}
}
