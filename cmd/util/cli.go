package util

import (
	"fmt"
	"github.com/camerondurham/ch/cmd/streams"
	"github.com/docker/docker/api/types/container"
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

type ContainerOpts struct {
	BuildOpts  *BuildOpts
	PullOpts   *PullOpts
	HostConfig *container.HostConfig

	//Volume      []string // maps to HostConfig.Volumes, consolidate and just use HostConfig
	//CapAdd      []string // HostConfig.CapAdd
	//SecurityOpt []string // HostConfig.SecurityOpt
	//Privileged  bool	 // HostConfig.Privileged

	Shell string
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
	Running() (map[string]string, error)
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

func SetEnvs(envs map[string]*ContainerOpts) {
	viper.Set("envs", envs)
}

func (cli *Cli) Running() (running map[string]string, err error) {
	if !viper.IsSet("running") {
		return nil, ErrDoesNotExist
	} else {
		running = make(map[string]string)
		err = viper.UnmarshalKey("running", &running)
		return
	}
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
