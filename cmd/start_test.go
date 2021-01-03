package cmd

import (
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
