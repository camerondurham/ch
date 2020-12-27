package util

import (
	"context"
	"errors"
	"github.com/camerondurham/ch/cmd/util/mocks"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestDockerService_CreateContainer(t *testing.T) {
	type fields struct {
		client DockerClient
		cc     *client.Client
	}
	type args struct {
		ctx           context.Context
		config        *container.Config
		containerName string
		hostConfig    *container.HostConfig
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockDockerClient(ctrl)
	d := NewDockerServiceFromClient(m)

	fails := "build-fails"
	success := "build-success"

	m.EXPECT().
		ContainerCreate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(fails)).
		Return(container.ContainerCreateCreatedBody{}, errors.New("cannot create body")).
		AnyTimes()

	m.EXPECT().
		ContainerCreate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(success)).
		Return(container.ContainerCreateCreatedBody{}, nil).
		AnyTimes()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    container.ContainerCreateCreatedBody
		wantErr bool
	}{
		{
			name: "Test Build Failure",
			fields: fields{
				client: d,
			},
			args: args{
				containerName: fails,
			},
			want:    container.ContainerCreateCreatedBody{},
			wantErr: true,
		},
		{
			name: "Test Build Success",
			fields: fields{
				client: d,
			},
			args: args{
				containerName: success,
			},
			want:    container.ContainerCreateCreatedBody{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DockerService{
				client: tt.fields.client,
				cc:     tt.fields.cc,
			}
			got, err := d.CreateContainer(tt.args.ctx, tt.args.config, tt.args.containerName, tt.args.hostConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateContainer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateContainer() got = %v, want %v", got, tt.want)
			}
		})
	}
}
