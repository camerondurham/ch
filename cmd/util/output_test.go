package util

import (
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
						Context:        ".",
						Tag:            "alpine-cursed-test2",
					},
					HostConfig: &HostConfig{
						Binds:       []string{"C:\\Users\\Cameron\\Projects\\ch\\tests\\cursed folder:/var"},
						SecurityOpt: []string{"seccomp:unconfined"},
						CapAdd:      []string{"SYS_PTRACE"},
					},
				},
			},
			want: "Name:\talpine-cursed-test2\n\tDockerfile:\ttests\\Dockerfile.alpine\n\tContext:\t.\n\t    Tag:\talpine-cursed-test2\n\tVolume:\tC:\\Users\\Cameron\\Projects\\ch\\tests\\cursed folder:/var\n\tSecOpt:\tseccomp:unconfined\n\tCapAdd:\tSYS_PTRACE\n",
		},

		/*
		   Name:	alpine-cursed-test2
		   	Dockerfile:	tests\Dockerfile.alpine
		   	Context:	.
		   	    Tag:	alpine-cursed-test2
		   	Volume:	C:\Users\Cameron\Projects\ch\tests\cursed folder:/var
		   	SecOpt:	seccomp:unconfined
		   	CapAdd:	SYS_PTRACE
		*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildString(tt.args.envName, tt.args.opts); got != tt.want {
				t.Errorf("buildString() = \n\"%v\"\n want \n\"%v\"\n", got, tt.want)
			}
		})
	}
}
