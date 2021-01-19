package cmd

import (
	"fmt"
	"github.com/camerondurham/ch/cmd/util"
	"github.com/camerondurham/ch/cmd/util/mocks"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func Test_removeEnvironment(t *testing.T) {
	type args struct {
		envName string
		envs    map[string]*util.ContainerOpts
		cli     *util.Cli
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockDockerClient(ctrl)
	m.EXPECT().ContainerList(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("no containers running")).AnyTimes()
	d := util.NewDockerServiceFromClient(m)

	cli := util.NewCliClientWithDockerService(m, d)

	tests := []struct {
		name    string
		args    args
		want    map[string]*util.ContainerOpts
		wantErr bool
	}{
		{
			name: "Basic Remove",
			args: args{
				envName: "test",
				envs: map[string]*util.ContainerOpts{
					"test": &util.ContainerOpts{},
				},
				cli: cli,
			},
			want:    map[string]*util.ContainerOpts{},
			wantErr: false,
		},
		{
			name: "Remove With Other Keys",
			args: args{
				envName: "test1",
				envs: map[string]*util.ContainerOpts{
					"test1": &util.ContainerOpts{},
					"test2": &util.ContainerOpts{},
				},
				cli: cli,
			},
			want: map[string]*util.ContainerOpts{
				"test2": &util.ContainerOpts{},
			},
			wantErr: false,
		},
		{
			name: "Environment Not Found",
			args: args{
				envName: "test3",
				envs: map[string]*util.ContainerOpts{
					"test1": &util.ContainerOpts{},
					"test2": &util.ContainerOpts{},
				},
				cli: cli,
			},
			want: map[string]*util.ContainerOpts{
				"test1": &util.ContainerOpts{},
				"test2": &util.ContainerOpts{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := removeEnvironment(tt.args.envName, tt.args.envs, tt.args.cli)
			if (err != nil) != tt.wantErr {
				t.Errorf("removeEnvironment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeEnvironment() got = %v, want %v", got, tt.want)
			}
		})
	}
}
