package cmd

import (
	"github.com/camerondurham/ch/cmd/util"
	"reflect"
	"testing"
)

func Test_removeEnvironment(t *testing.T) {
	type args struct {
		envName string
		envs    map[string]*util.ContainerOpts
	}
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
			got, err := removeEnvironment(tt.args.envName, tt.args.envs)
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
