package cmd

import (
	"fmt"
	//"github.com/docker/cli/cli/streams"
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
	Volume    string
	Shell     string
}

// Streams is an interface which exposes the standard input and output streams
type Streams interface {
	In() *In
	Out() *Out
	Err() io.Writer
}

type Cli interface {
	Client() client.APIClient
	Out() *Out
	In() *In
	Err() io.Writer
}

type CliClient struct {
	in           *In
	out          *Out
	err          io.Writer
	dockerClient client.APIClient
}

func (cli *CliClient) Client() client.APIClient {
	return cli.dockerClient
}

func (cli *CliClient) In() *In {
	return cli.in
}

func (cli *CliClient) Out() *Out {
	return cli.out
}

func (cli *CliClient) Err() io.Writer {
	return cli.err
}

func NewCliClient() (*CliClient, error) {
	cliClient := &CliClient{}

	_, dockerClient := DockerClientInitOrDie()

	if dockerClient != nil {
		cliClient.dockerClient = dockerClient
	} else {
		return nil, fmt.Errorf("error creating docker client")
	}

	stdin, stdout, stderr := term.StdStreams()
	cliClient.in = NewIn(stdin)
	cliClient.out = NewOut(stdout)
	cliClient.err = stderr

	return cliClient, nil
}

/*
// DockerCli is an instance the docker command line client.
// Instances of the client can be returned from NewDockerCli.
type DockerCli struct {
	configFile         *configfile.ConfigFile
	in                 *streams.In
	out                *streams.Out
	err                io.Writer
	client             client.APIClient
	serverInfo         ServerInfo
	clientInfo         *ClientInfo
	contentTrust       bool
	contextStore       store.Store
	currentContext     string
	dockerEndpoint     docker.Endpoint
	contextStoreConfig store.Config
}

// Cli represents the docker command line client.
type Cli interface {
	Client() client.APIClient
	Out() *streams.Out
	Err() io.Writer
	In() *streams.In
	SetIn(in *streams.In)
	Apply(ops ...DockerCliOption) error
	ConfigFile() *configfile.ConfigFile
	ServerInfo() ServerInfo
	ClientInfo() ClientInfo
	NotaryClient(imgRefAndAuth trust.ImageRefAndAuth, actions []string) (notaryclient.Repository, error)
	DefaultVersion() string
	ManifestStore() manifeststore.Store
	RegistryClient(bool) registryclient.RegistryClient
	ContentTrustEnabled() bool
	ContextStore() store.Store
	CurrentContext() string
	StackOrchestrator(flagValue string) (Orchestrator, error)
	DockerEndpoint() docker.Endpoint
}
*/
