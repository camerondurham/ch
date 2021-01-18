package util

import (
	"github.com/docker/go-connections/nat"
	"testing"
)

func Test_buildString(t *testing.T) {
	type args struct {
		envName string
		opts    *ContainerOpts
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "Empty opts",
			args: args{
				envName: "test",
				opts: &ContainerOpts{
					BuildOpts:  nil,
					PullOpts:   nil,
					HostConfig: nil,
					Shell:      "",
				},
			},
			want: "Name:\ttest\n",
		},
		{
			// TODO: make this dynamic
			name: "Nonempty opts",
			args: args{
				envName: "alpine-cursed-test2",
				opts: &ContainerOpts{
					BuildOpts: &BuildOpts{
						DockerfilePath: "tests\\Dockerfile.alpine",
						Context:        "C:\\Users\\Cameron\\Projects\\ch\\scripts\\test-scripts",
						Tag:            "alpine-cursed-test2",
					},
					HostConfig: &HostConfig{
						Binds:       []string{"C:\\Users\\Cameron\\Projects\\ch\\tests\\cursed folder:/var"},
						SecurityOpt: []string{"seccomp:unconfined"},
						CapAdd:      []string{"SYS_PTRACE"},
					},
				},
			},
			want: "Name:\talpine-cursed-test2\n\tDockerfile:\ttests\\Dockerfile.alpine\n\tContext:\tC:\\Users\\Cameron\\Projects\\ch\\scripts\\test-scripts\n\t    Tag:\talpine-cursed-test2\n\tVolume:\tC:\\Users\\Cameron\\Projects\\ch\\tests\\cursed folder:/var\n\tSecOpt:\tseccomp:unconfined\n\tCapAdd:\tSYS_PTRACE\n",
		},
		{
			name: "Nonempty opts.PortBindings",
			args: args{
				envName: "csci350",
				opts: &ContainerOpts{
					PullOpts: &PullOpts{ImageName: "camerondurham/xv6-docker:latest"},
					HostConfig: &HostConfig{
						SecurityOpt: []string{"seccomp:unconfined"},
						CapAdd:      []string{"SYS_PTRACE"},
						PortBindings: map[nat.Port][]nat.PortBinding{
							nat.Port("22"): {nat.PortBinding{HostPort: "22"}}},
					},
				},
			},
			want: "Name:\tcsci350\n\tImage:\tcamerondurham/xv6-docker:latest\n\tSecOpt:\tseccomp:unconfined\n\tCapAdd:\tSYS_PTRACE\n\tPort:\t22:22\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildString(tt.args.envName, tt.args.opts); got != tt.want {
				t.Errorf("buildString() = \n\"%v\"\n want \n\"%v\"\n", got, tt.want)
			}
		})
	}
}
