package cmd

import (
	"github.com/camerondurham/ch/cmd/util"
	"github.com/docker/docker/api/types/container"
	"reflect"
	"testing"
)

func Test_updateRunning(t *testing.T) {
	type args struct {
		running     map[string]string
		containerID string
		envName     string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "Empty Running Map",
			args: args{
				running:     nil,
				containerID: "1234567890",
				envName:     "test",
			},
			want: map[string]string{
				"test": "1234567890",
			},
		},
		{
			name: "Non-Empty Running Map",
			args: args{
				running: map[string]string{
					"existing": "1234567890",
				},
				containerID: "1234567890",
				envName:     "test",
			},
			want: map[string]string{
				"existing": "1234567890",
				"test":     "1234567890",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := updateRunning(tt.args.running, tt.args.containerID, tt.args.envName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("updateRunning() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createHostConfig(t *testing.T) {
	type args struct {
		containerOpts *util.ContainerOpts
	}

	tests := []struct {
		name string
		args args
		want *container.HostConfig
	}{
		{
			name: "Empty ContainerOpts",
			args: args{
				containerOpts: &util.ContainerOpts{},
			},
			want: &container.HostConfig{},
		},
		{
			name: "Has HostConfig",
			args: args{
				containerOpts: &util.ContainerOpts{
					HostConfig: &util.HostConfig{
						Binds:       []string{"/home/cam/projects:/work"},
						SecurityOpt: []string{"seccomp:unconfined"},
						Privileged:  true,
						CapAdd:      []string{"SYS_PTRACE"},
					},
				},
			},
			want: &container.HostConfig{
				Binds:       []string{"/home/cam/projects:/work"},
				Privileged:  true,
				SecurityOpt: []string{"seccomp:unconfined"},
				CapAdd:      []string{"SYS_PTRACE"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createHostConfig(tt.args.containerOpts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createHostConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
