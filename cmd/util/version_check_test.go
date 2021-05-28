package util

import (
	"fmt"
	"testing"
)

func TestLatestVersion(t *testing.T) {
	type args struct {
		repository string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Nominal: camerondurham/ch",
			args: args{repository: "camerondurham/ch"},
			want: "https://api.github.com/repos/camerondurham/ch/releases/latest",
		},
		{
			name: "Nominal: ph3rin/sacc",
			args: args{repository: "ph3rin/sacc"},
			want: "https://api.github.com/repos/ph3rin/sacc/releases/latest",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LatestVersion(tt.args.repository); got != tt.want {
				t.Errorf("LatestVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLatestVersion(t *testing.T) {
	type args struct {
		getRequest callback
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Nominal",
			args: args{
				getRequest: func(s string) (map[string]interface{}, error) {
					return map[string]interface{}{
						"tag_name": "v0.2.2",
					}, nil
				},
			},
			want:    "v0.2.2",
			wantErr: false,
		},
		{
			name: "Error",
			args: args{
				getRequest: func(s string) (map[string]interface{}, error) {
					return nil, fmt.Errorf("error making GET request")
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLatestVersion(tt.args.getRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLatestVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetLatestVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}
