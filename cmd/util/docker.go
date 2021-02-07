package util

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/camerondurham/ch/cmd/streams"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/pkg/jsonmessage"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	specs "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/docker/docker/api/types"
)

type DockerClient interface {
	ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error)
	ImagePull(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error)
	ImageBuild(ctx context.Context, buildContext io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error)
	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *specs.Platform, containerName string) (container.ContainerCreateCreatedBody, error)
	ContainerRemove(ctx context.Context, containerID string, options types.ContainerRemoveOptions) error
	ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error
	ContainerStop(ctx context.Context, containerID string, timeout *time.Duration) error
}

type DockerService struct {
	client    DockerClient
	cc        *client.Client
	validator *Validator
}

// Docker Build Response
type ErrorDetail struct {
	Code    int    `json:",string"`
	Message string `json:"message"`
}

func NewDockerService() (*DockerService, error) {
	_, cli, err := dockerClientInit()
	if err != nil {
		fmt.Printf("error creating docker client: %v", err)
		return nil, errors.New("cannot create Docker client")
	} else {
		return &DockerService{client: cli, cc: cli, validator: &Validator{}}, nil
	}
}

func NewDockerServiceFromClient(cli DockerClient) *DockerService {
	return &DockerService{client: cli}
}

// GetRunning lists running containers like docker ps
func (d *DockerService) GetRunning(filters filters.Args, print bool) ([]types.Container, error) {
	containers, err := d.ContainerList(context.Background(), types.ContainerListOptions{Filters: filters})
	if err != nil {
		DebugPrint(fmt.Sprintf("error listing docker containers: %v", err))
		return nil, fmt.Errorf("failed to get running containers\nplease verify that Docker is running")
	}

	if print {
		for _, c := range containers {
			fmt.Printf("%s %s\n", c.ID[:10], c.Image)
		}
	}
	return containers, nil
}

// PullImage downloads a Docker image from Docker Hub
func (d *DockerService) PullImage(ctx context.Context, out *streams.Out, imageName string) error {
	imagePullOutput, err := d.ImagePull(ctx, imageName, types.ImagePullOptions{})

	if err != nil {
		return err
	}

	type ImagePullOutput struct {
		Status         string `json:"status"`
		Error          string `json:"error"`
		Progress       string `json:"progress"`
		ProgressDetail struct {
			Current int `json:"current"`
			Total   int `json:"total"`
		} `json:"progressDetail"`
	}

	defer imagePullOutput.Close()

	return jsonmessage.DisplayJSONMessagesToStream(imagePullOutput, out, nil)
}

// BuildImageWithContext accepts a build context path and relative Dockerfile path
func (d *DockerService) BuildImageWithContext(ctx context.Context, dockerfile string, contextDirPath string, imageTagName string) (err error) {
	contextPath := contextDirPath
	contextPath = d.validator.GetAbs(contextPath)

	if !d.validator.ValidPath(contextPath) {
		return errors.New(fmt.Sprintf("context path does not exist: %v", contextPath))
	}

	// TODO: use ioutil.TempFile
	dir, err := ioutil.TempDir("", "docker-context")
	if err != nil {
		return fmt.Errorf("error creating TempDir: %v", err)
	}
	defer os.RemoveAll(dir)
	contextTarball := filepath.Join(dir, "context.tar")

	DebugPrint(fmt.Sprintf("dockerfile context file: %s\n", contextPath))
	DebugPrint(fmt.Sprintf("output filename: %s\n", contextTarball))

	contextTarReader, err := contextReader(contextPath)
	if err != nil {
		return err
	}

	buildResponse, err := d.ImageBuild(ctx, contextTarReader, types.ImageBuildOptions{
		Context:    contextTarReader,
		Tags:       []string{imageTagName},
		Dockerfile: dockerfile,
		Remove:     true,
	})

	if err != nil {
		DebugPrint(fmt.Sprintf("unable to build docker image: \n%v\n", err))
		return err
	}

	type BuildOutput struct {
		Stream      string       `json:"stream"`
		ErrorDetail *ErrorDetail `json:"errorDetail"`
		Error       string       `json:"error,omitempty"`
	}

	defer buildResponse.Body.Close()

	DebugPrint(buildResponse.OSType)

	rd := bufio.NewReader(buildResponse.Body)

	var retErr error

	for {

		// there must be a better way than parsing the output to figure out if a build failed??
		str, err := rd.ReadString('\n')
		if err == io.EOF {
			retErr = nil
			break
		} else {
			var msg BuildOutput

			err = json.Unmarshal([]byte(str), &msg)

			if err != nil {
				DebugPrint(fmt.Sprintf("error unmarshalling str: [%s] \n error: %v", str, err))
			}

			if msg.Error != "" {
				retErr = fmt.Errorf("error building image:\n%v", msg.ErrorDetail.Message)
				break
			}
		}
	}
	return retErr
}

// CreateContainer create container with name
func (d *DockerService) CreateContainer(ctx context.Context, config *container.Config, containerName string, hostConfig *container.HostConfig) (container.ContainerCreateCreatedBody, error) {

	createdContainerResp, err := d.ContainerCreate(ctx, config, hostConfig, nil, nil, containerName)

	if err != nil {
		DebugPrint(fmt.Sprint("unable to create container: ", err))
		return createdContainerResp, err
	}

	return createdContainerResp, nil
}

// RemoveContainer delete container
func (d *DockerService) RemoveContainer(ctx context.Context, containerName string) error {

	DebugPrint(fmt.Sprintf("removing container[%s]", containerName))

	if err := d.ContainerRemove(context.Background(), containerName, types.ContainerRemoveOptions{}); err != nil {
		DebugPrint(fmt.Sprint("error removing container: ", err))
		return err
	} else {
		return nil
	}
}

// StartContainer with given name
func (d *DockerService) StartContainer(ctx context.Context, containerID string) error {

	DebugPrint(fmt.Sprintf("starting container [%s]", containerID))

	if err := d.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		DebugPrint(fmt.Sprint("error starting container: ", err))
		return err
	} else {
		return nil
	}
}

// StopContainer from running
func (d *DockerService) StopContainer(ctx context.Context, containerID string, timeout *time.Duration) error {

	DebugPrint(fmt.Sprintf("removing container [%s]...", containerID))

	if err := d.ContainerStop(ctx, containerID, nil); err != nil {
		DebugPrint(fmt.Sprint("error stopping container: ", err))
		return err
	} else {
		return nil
	}
}

// Helper functions

func dockerClientInit() (ctx context.Context, cli *client.Client, err error) {
	ctx = context.Background()
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		DebugPrint(fmt.Sprintf("error creating Docker client: %v\nare you sure the client is running?\n", err))
		return nil, nil, err
	}
	return ctx, cli, nil
}

// contextReader reads path in to tar buffer
func contextReader(contextPath string) (contextReader *bytes.Reader, err error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer func() {
		err = tw.Close()
	}()

	path := filepath.Clean(contextPath)

	walker := func(file string, finfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// fill in header info using func FileInfoHeader
		hdr, err := tar.FileInfoHeader(finfo, finfo.Name())
		if err != nil {
			return err
		}
		relFilePath := file
		if filepath.IsAbs(path) {
			relFilePath, err = filepath.Rel(path, file)
			if err != nil {
				return err
			}
		}

		// ensure header has relative file path
		hdr.Name = relFilePath

		if err2 := tw.WriteHeader(hdr); err2 != nil {
			return err2
		}

		// if file is a dir, don't continue
		if finfo.Mode().IsDir() {
			return nil
		}

		// add file to tar
		srcFile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer srcFile.Close()
		_, err = io.Copy(tw, srcFile)
		if err != nil {
			return err
		}
		return nil
	}

	// build tar
	if err := filepath.Walk(path, walker); err != nil {
		fmt.Printf("failed to add %s to tar buffer: %s\n", path, err)
	}

	contextTarReader := bytes.NewReader(buf.Bytes())

	return contextTarReader, nil
}

// Interface Implementation

func (d *DockerService) ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error) {
	return d.client.ContainerList(ctx, options)
}

func (d *DockerService) ImagePull(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error) {
	return d.client.ImagePull(ctx, refStr, options)
}

func (d *DockerService) ImageBuild(ctx context.Context, buildContext io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	return d.client.ImageBuild(ctx, buildContext, options)
}

func (d *DockerService) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *specs.Platform, containerName string) (container.ContainerCreateCreatedBody, error) {
	return d.client.ContainerCreate(ctx, config, hostConfig, networkingConfig, platform, containerName)
}

func (d *DockerService) ContainerRemove(ctx context.Context, containerID string, options types.ContainerRemoveOptions) error {
	return d.client.ContainerRemove(ctx, containerID, options)
}

func (d *DockerService) ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error {
	return d.client.ContainerStart(ctx, containerID, options)
}

func (d *DockerService) ContainerStop(ctx context.Context, containerID string, timeout *time.Duration) error {
	return d.client.ContainerStop(ctx, containerID, timeout)
}
