package cmd

// TODO: rewrite as implementations on the cli

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/cli/cli"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	"github.com/docker/docker/api/types"
)

// Docker Build Response

type ErrorDetail struct {
	Code    int    `json:",string"`
	Message string `json:"message"`
}

type BuildOutput struct {
	Stream      string       `json:"stream"`
	ErrorDetail *ErrorDetail `json:"errorDetail"`
	Error       string       `json:"error,omitempty"`
}

func DockerClientInitOrDie() (ctx context.Context, cli *client.Client) {
	ctx = context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Printf("error creating Docker client: %v\nare you sure the client is running?\n", err)
		os.Exit(1)
	}
	return ctx, cli
}

// ListRunning lists running containers like docker ps
func ListRunning(cli *client.Client) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range containers {
		fmt.Printf("%s %s\n", c.ID[:10], c.Image)
	}
}

// ContextReader reads path in to tar buffer
func ContextReader(contextPath string) (contextReader *bytes.Reader, err error) {
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

// PullImage downloads a Docker image from Docker Hub
func PullImage(ctx context.Context, cli *client.Client, imageName string) error {
	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})

	if err != nil {
		return err
	}

	_, err = io.Copy(os.Stdout, out)
	if err != nil {
		return err
	}

	return nil
}

// BuildImageWithContext accepts a build context path and relative Dockerfile path
func BuildImageWithContext(ctx context.Context, cli *client.Client, dockerfile string, contextDirPath string, imageTagName string) (err error) {
	contextPath, err := filepath.Abs(contextDirPath)
	if err != nil {
		log.Printf("error finding abs path: %v", err)
		return err
	}
	contextTarball := fmt.Sprintf("/tmp/%s.tar", filepath.Base(contextPath))

	DebugPrint(fmt.Sprintf("dockerfile context file: %s\n", contextPath))
	DebugPrint(fmt.Sprintf("output filename: %s\n", contextTarball))

	contextTarReader, err := ContextReader(contextPath)
	if err != nil {
		return err
	}

	buildResponse, err := cli.ImageBuild(ctx, contextTarReader, types.ImageBuildOptions{
		Context:    contextTarReader,
		Tags:       []string{imageTagName},
		Dockerfile: dockerfile,
		Remove:     true,
	})

	if err != nil {
		log.Printf("unable to build docker image: %v", err)
		return err
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
func CreateContainer(ctx context.Context, cli *client.Client, config *container.Config, containerName string) container.ContainerCreateCreatedBody {
	createdContainerResp, err := cli.ContainerCreate(ctx, config, nil, nil, nil, containerName)

	if err != nil {
		log.Fatal("unable to create container: ", err)
	}

	return createdContainerResp
}

// RemoveContainer delete container
func RemoveContainer(ctx context.Context, cli *client.Client, containerName string) {
	DebugPrint(fmt.Sprintf("removing container[%s]...", containerName))
	if err := cli.ContainerRemove(context.Background(), containerName, types.ContainerRemoveOptions{}); err != nil {
		log.Fatal("error removing container: ", err)
	}
}

// StartContainer with given name
func StartContainer(ctx context.Context, cli *client.Client, containerID string) {
	if err := cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		log.Fatal("unable to start container: ", err)
	}
}

// StopContainer from running
func StopContainer(ctx context.Context, cli *client.Client, containerID string, timeout *time.Duration) {

	DebugPrint(fmt.Sprintf("removing container [%s]...", containerID))

	if err := cli.ContainerStop(ctx, containerID, nil); err != nil {
		log.Fatal("error stopping container: ", err)
	}
}

// CreateExecInteractive creates an exec config to run an exec process
func CreateExecInteractive(ctx context.Context, dockerClient *client.Client, cliClient Cli, container string, config types.ExecConfig) error {
	if _, err := dockerClient.ContainerInspect(ctx, container); err != nil {
		return err
	}

	// avoid config Detach check if tty is correct

	response, err := dockerClient.ContainerExecCreate(ctx, container, config)
	if err != nil {
		return err
	}
	execID := response.ID
	if execID == "" {
		return errors.New("exec ID empty")
	}

	if config.Detach {
		execStartCheck := types.ExecStartCheck{
			Detach: config.Tty,
			Tty:    config.Tty,
		}
		return dockerClient.ContainerExecStart(ctx, execID, execStartCheck)
	}
	return interactiveExec(ctx, dockerClient, cliClient, &config, execID)

}

func interactiveExec(ctx context.Context, dockerClient *client.Client, cliClient Cli, execConfig *types.ExecConfig, execID string) error {
	var (
		out, stderr io.Writer
		in          io.ReadCloser
	)

	// attach stdin, possibly add more functionality later
	in = os.Stdin
	out = os.Stdout

	// attach to os.Stderr only if not tty?
	stderr = os.Stdout

	resp, err := dockerClient.ContainerExecAttach(ctx, execID, types.ExecStartCheck{Tty: true})

	if err != nil {
		log.Fatal("error attaching exec to container: ", err)
	}
	defer resp.Close()

	errCh := make(chan error, 1)

	go func() {
		defer close(errCh)
		errCh <- func() error {

			// get streamer as hijackedIOStreamer
			streamer := hijackedIOStreamer{
				streams:      cliClient,
				inputStream:  in,
				outputStream: out,
				errorStream:  stderr,
				resp:         resp,
				tty:          true,
			}

			// return streamer.stream(ctx)
			return streamer.stream(ctx)
		}()
	}()

	// ignore check if config wants a terminal and has appropriate Tty size for now

	// check MonitorTtySize
	if err := <-errCh; err != nil {
		// TODO: debug print
		log.Printf("Error hijack: %v", err)
		return err
	}

	return getExecExitStatus(ctx, dockerClient, execID)
}
func getExecExitStatus(ctx context.Context, dockerClient client.ContainerAPIClient, execID string) error {
	resp, err := dockerClient.ContainerExecInspect(ctx, execID)
	if err != nil {
		// daemon probably died
		if !client.IsErrConnectionFailed(err) {
			return err
		}
		return cli.StatusError{StatusCode: -1}
	}
	status := resp.ExitCode
	if status != 0 {
		return cli.StatusError{StatusCode: status}
	}
	return nil
}
