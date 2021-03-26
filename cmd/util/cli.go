package util

import (
	"errors"
	"fmt"
	"github.com/camerondurham/ch/cmd/streams"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
	"github.com/spf13/viper"
	"os"

	"github.com/moby/term"
	"io"
)

type BuildOpts struct {
	DockerfilePath string
	Context        string
	Tag            string
}

type PullOpts struct {
	ImageName string
}

type HostConfig struct {
	Binds        []string          // List of volume bindings for this container
	SecurityOpt  []string          // List of string values to customize labels for MLS systems, such as SELinux.
	PortBindings nat.PortMap       // Map of host to container ports
	Privileged   bool              // Is the container in privileged mode
	CapAdd       strslice.StrSlice // List of kernel capabilities to add to the container
}

type ContainerOpts struct {
	BuildOpts  *BuildOpts
	PullOpts   *PullOpts
	HostConfig *HostConfig
	Shell      string
}

// TODO: add config/env settings to use Cli in other commands
type Cli struct {
	in              *streams.In
	out             *streams.Out
	err             io.Writer
	dockerAPIClient *DockerAPIService
	dockerService   *DockerService
	validator       *Validator
}

type ContainerClient interface {
	ApiClient() *DockerAPIService
	DockerClient() *DockerService
	Out() *streams.Out
	In() *streams.In
	Err() io.Writer
	Containers() map[string]*ContainerOpts
	Container(envName string) (*types.Container, error)
	ContainerIsRunning(envName string) bool
	Validator() *Validator
}

func (cli *Cli) DockerClient() *DockerService {
	return cli.dockerService
}

func (cli *Cli) ApiClient() *DockerAPIService {
	return cli.dockerAPIClient
}

func (cli *Cli) In() *streams.In {
	return cli.in
}

func (cli *Cli) Out() *streams.Out {
	return cli.out
}

func (cli *Cli) Err() io.Writer {
	return cli.err
}

func (cli *Cli) Containers() map[string]*ContainerOpts {
	envs, err := GetEnvs()
	if err != nil {
		fmt.Printf("error retrieving envs: %v\n", err)
		os.Exit(1)
	}
	return envs
}

func GetEnvs() (envs map[string]*ContainerOpts, err error) {
	if !viper.IsSet("envs") {
		return nil, ErrDoesNotExist
	}

	envs = make(map[string]*ContainerOpts)
	err = viper.UnmarshalKey("envs", &envs)
	return
}

func PrintEnvNotFoundMsg(envName string) {
	fmt.Printf("no such environment: %v\n", envName)
	envs, err := GetEnvs()
	if err == nil {
		fmt.Printf("\nAvailable environments:\n\n")
		for k, _ := range envs {
			fmt.Printf("\t%s\n", k)
		}
	}
}

func SetEnvs(envs map[string]*ContainerOpts) {
	viper.Set("envs", envs)
}

func (cli *Cli) Container(envName string) (*types.Container, error) {
	f := filters.NewArgs(filters.Arg("name", envName))
	c, err := cli.DockerClient().GetRunning(f, false)

	if err != nil {
		return nil, errors.New("failed to get running environments")
	}

	if len(c) != 1 {
		return nil, errors.New("environment not running")
	}
	return &c[0], nil
}

func (cli *Cli) ContainerIsRunning(envName string) bool {
	c, err := cli.Container(envName)
	return err == nil && c != nil
}

func (cli *Cli) Validator() *Validator {
	return cli.validator
}

func NewCliClient() (*Cli, error) {
	cliClient := &Cli{}

	dockerService, err := NewDockerService()

	if err == nil {
		cliClient.dockerAPIClient = NewDockerAPIService(*dockerService)
		cliClient.dockerService = dockerService
	} else {
		return nil, fmt.Errorf("error creating docker client")
	}

	stdin, stdout, stderr := term.StdStreams()
	cliClient.in = streams.NewIn(stdin)
	cliClient.out = streams.NewOut(stdout)
	cliClient.err = stderr
	cliClient.validator = &Validator{}

	return cliClient, nil
}

func NewCliClientWithDockerService(dockerClient DockerClient, dockerService *DockerService) *Cli {
	cliClient := &Cli{}

	cliClient.dockerAPIClient = NewDockerAPIService(*dockerService)
	cliClient.dockerService = dockerService

	stdin, stdout, stderr := term.StdStreams()
	cliClient.in = streams.NewIn(stdin)
	cliClient.out = streams.NewOut(stdout)
	cliClient.err = stderr
	cliClient.validator = &Validator{}

	return cliClient
}
