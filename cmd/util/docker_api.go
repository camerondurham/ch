package util

import (
	"context"
	"errors"
	"fmt"
	"github.com/camerondurham/ch/cmd/streams"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"log"
)

type DockerAPI interface {
	CreateExecInteractive(ctx context.Context, cliClient ContainerClient, container string, config types.ExecConfig) error

	containerInspect(ctx context.Context, container string) (types.ContainerJSON, error)
	containerExecCreate(ctx context.Context, container string, config types.ExecConfig) (types.IDResponse, error)
	containerExecStart(ctx context.Context, execID string, config types.ExecStartCheck) error
	containerExecAttach(ctx context.Context, execID string, config types.ExecStartCheck) (types.HijackedResponse, error)
	interactiveExec(ctx context.Context, cliClient ContainerClient, execConfig *types.ExecConfig, execID string) error
}

type DockerAPIService struct {
	client client.APIClient
}

// TODO: improve this interface
func NewDockerAPIService(client *client.Client) *DockerAPIService {
	return &DockerAPIService{client: client}
}

func (d *DockerAPIService) containerInspect(ctx context.Context, container string) (types.ContainerJSON, error) {
	return d.client.ContainerInspect(ctx, container)
}
func (d *DockerAPIService) containerExecCreate(ctx context.Context, container string, config types.ExecConfig) (types.IDResponse, error) {
	return d.client.ContainerExecCreate(ctx, container, config)
}
func (d *DockerAPIService) containerExecStart(ctx context.Context, execID string, config types.ExecStartCheck) error {
	return d.client.ContainerExecStart(ctx, execID, config)
}
func (d *DockerAPIService) containerExecAttach(ctx context.Context, execID string, config types.ExecStartCheck) (types.HijackedResponse, error) {
	return d.client.ContainerExecAttach(ctx, execID, config)
}

// CreateExecInteractive creates an exec config to run an exec process
func (d *DockerAPIService) CreateExecInteractive(ctx context.Context, cliClient ContainerClient, container string, config types.ExecConfig) error {
	if _, err := d.containerInspect(ctx, container); err != nil {
		return err
	}

	// avoid config Detach check if tty is correct

	response, err := d.containerExecCreate(ctx, container, config)
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
		return d.containerExecStart(ctx, execID, execStartCheck)
	}
	return d.interactiveExec(ctx, cliClient, &config, execID)

}

func (d *DockerAPIService) interactiveExec(ctx context.Context, cliClient ContainerClient, execConfig *types.ExecConfig, execID string) error {
	var (
		out, stderr io.Writer
		in          io.ReadCloser
	)

	// attach stdin, possibly add more functionality later
	in = cliClient.In()
	out = cliClient.Out()

	// attach to os.Stderr only if not tty?
	stderr = cliClient.Err()

	resp, err := d.containerExecAttach(ctx, execID, types.ExecStartCheck{Tty: true})

	if err != nil {
		log.Fatal("error attaching exec to container: ", err)
	}
	defer resp.Close()

	errCh := make(chan error, 1)

	go func() {
		defer close(errCh)
		errCh <- func() error {

			// get streamer as hijackedIOStreamer
			streamer := streams.HijackedIOStreamer{
				Streams:      cliClient,
				InputStream:  in,
				OutputStream: out,
				ErrorStream:  stderr,
				Resp:         resp,
				Tty:          execConfig.Tty,
			}

			return streamer.Stream(ctx)
		}()
	}()

	// ignore check if config wants a terminal and has appropriate Tty size for now

	// check MonitorTtySize
	if err := <-errCh; err != nil {
		DebugPrint(fmt.Sprintf("Error hijack: %v", err))
		return err
	}

	return getExecExitStatus(ctx, d.client, execID)
}

func getExecExitStatus(ctx context.Context, dockerClient client.ContainerAPIClient, execID string) error {
	resp, err := dockerClient.ContainerExecInspect(ctx, execID)
	if err != nil {
		// daemon probably died
		if !client.IsErrConnectionFailed(err) {
			return err
		}
		return errors.New(fmt.Sprintf("error status code: %v,\nmessage: %v ", -1, err))
	}
	status := resp.ExitCode
	if status != 0 {
		return errors.New(fmt.Sprintf("error status code: %v,\nmessage: %v ", status, err))
	}
	return nil
}
