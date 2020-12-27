package util

import (
	"context"
	"errors"
	"github.com/camerondurham/ch/cmd/util/mocks"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/golang/mock/gomock"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"testing"
)

func TestDockerService_CreateContainer(t *testing.T) {
	//type fields struct {
	//	client *client.Client
	//}
	//type args struct {
	//	ctx           context.Context
	//	config        *container.Config
	//	containerName string
	//	hostConfig    *container.HostConfig
	//}
	//
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockDockerClient(ctrl)
	m.EXPECT().
		ContainerCreate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), "fail-build").
		DoAndReturn(func(_ context.Context, _ *container.Config, _ *container.HostConfig, _ *network.NetworkingConfig, _ *specs.Platform, containerName string) (container.ContainerCreateCreatedBody, error) {
			return container.ContainerCreateCreatedBody{}, errors.New("cannot create body")
		})

	name := "Create Container Nominal"

	t.Run(name, func(t *testing.T) {
		d := NewDockerServiceFromClient(m)
		if _, err := d.CreateContainer(context.Background(), &container.Config{}, "fail-build", &container.HostConfig{}); err == nil {
			t.Errorf("expected error, none returned")
		}
	})
}
