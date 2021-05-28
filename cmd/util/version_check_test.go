package util

import "testing"

func TestGetLatestVersion(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		// TODO: remove this test after API gets bumped. I'm too lazy to handle mocking right now.
		{
			name:    "Nominal: basic testing",
			want:    "v0.2.2",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLatestVersion()
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
