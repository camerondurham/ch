package util

import (
	"fmt"
	"github.com/camerondurham/ch/cmd/streams"
	"github.com/spf13/viper"
	"os"

	"github.com/docker/docker/client"
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
	BuildOpts *BuildOpts
	PullOpts  *PullOpts
	Volume    []string
	Shell     string
}

// TODO: add config/env settings to use Cli in other commands
type Cli struct {
	in              *streams.In
	out             *streams.Out
	err             io.Writer
	dockerClient    *client.Client
	dockerAPIClient client.APIClient
}

type ContainerClient interface {
	ApiClient() client.APIClient
	Client() *client.Client
	Out() *streams.Out
	In() *streams.In
	Err() io.Writer
	Containers() map[string]*ContainerOpts
	Running() (map[string]string, error)
}

func (cli *Cli) ApiClient() client.APIClient {
	return cli.dockerAPIClient
}

func (cli *Cli) Client() *client.Client {
	return cli.dockerClient
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

func NewCliClient() (*Cli, error) {
	cliClient := &Cli{}

	_, dockerClient := DockerClientInitOrDie()

	if dockerClient != nil {
		cliClient.dockerAPIClient = dockerClient
		cliClient.dockerClient = dockerClient
	} else {
		return nil, fmt.Errorf("error creating docker client")
	}

	stdin, stdout, stderr := term.StdStreams()
	cliClient.in = streams.NewIn(stdin)
	cliClient.out = streams.NewOut(stdout)
	cliClient.err = stderr

	return cliClient, nil
}
