package dockerutil

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/camerondurham/ch/cmd/util"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	"github.com/docker/docker/api/types"
)

// ListRunning lists running containers like docker ps
func ListRunning(cli *client.Client) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
}

// ContextReader reads path in to tar buffer
func ContextReader(contextPath string) (contextReader *bytes.Reader, err error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	path := filepath.Clean(contextPath)

	if err != nil {
		fmt.Println(err)
	}

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

		if err := tw.WriteHeader(hdr); err != nil {
			return err
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

// BuildImageWithContext accepts a build context path and relative Dockerfile path
func BuildImageWithContext(ctx context.Context, cli *client.Client, dockerfile string, contextDirPath string, imageTagName string) error {
	contextPath, err := filepath.Abs(contextDirPath)
	if err != nil {
		log.Printf("error finding abs path: %v", err)
		return err
	}
	contextTarball := fmt.Sprintf("/tmp/%s.tar", filepath.Base(contextPath))

	util.DebugPrint(fmt.Sprintf("dockerfile context file: %s\n", contextPath))
	util.DebugPrint(fmt.Sprintf("output filename: %s\n", contextTarball))

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

	util.DebugPrint(buildResponse.OSType)

	_, err = io.Copy(os.Stdout, buildResponse.Body)
	if err != nil {
		log.Fatal("unable to read image build response: ", err)
		return err
	}

	return nil
}

// CreateContainer create container with name
func CreateContainer(ctx context.Context, cli *client.Client, config *container.Config, containerName string) container.ContainerCreateCreatedBody {
	createdContainerResp, err := cli.ContainerCreate(ctx, config, nil, nil, containerName)

	if err != nil {
		log.Fatal("unable to create container: ", err)
	}

	return createdContainerResp
}

// RemoveContainer delete container
func RemoveContainer(ctx context.Context, cli *client.Client, containerName string) {
	util.DebugPrint(fmt.Sprintf("removing container[%s]...", containerName))
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

	util.DebugPrint(fmt.Sprintf("removing container [%s]...", containerID))

	if err := cli.ContainerStop(ctx, containerID, nil); err != nil {
		log.Fatal("error stopping container: ", err)
	}
}
