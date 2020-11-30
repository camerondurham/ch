package cmd

import (
	"bytes"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"reflect"
	"testing"
)

//func TestCreate(t *testing.T) {
//	cmd := rootCmd
//
//	args := []string{"create", "test1"}
//	cmd.SetArgs(args)
//	b := bytes.NewBufferString("")
//	cmd.SetOut(b)
//	cmd.Execute()
//	out, err := ioutil.ReadAll(b)
//	if err != nil {
//		t.Fatal(err)
//	}
//	expectsContains := "requires at least 1 arg"
//	if strings.Contains(string(out), expectsContains) == false {
//		t.Fatalf("expected to contain: %s, got: %s", expectsContains, string(out))
//	} else {
//		t.Log("correct error")
//	}
//
//}

func TestCreateCmd(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name string
		args args
	}{
		{"create",
			args{createCmd, []string{"create", "simNominal"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBufferString("")
			tt.args.cmd.SetOut(b)
			tt.args.cmd.Run(tt.args.cmd, tt.args.args)
			out, err := ioutil.ReadAll(b)
			if err != nil {
				t.Fatal(err)
			}
			log.Printf("output: %s", out)
		})
	}
}

func Test_parseContainerOpts(t *testing.T) {
	type args struct {
		cmd             *cobra.Command
		environmentName string
	}
	tests := []struct {
		name    string
		args    args
		want    *ContainerOpts
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseContainerOpts(tt.args.cmd, tt.args.environmentName)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseContainerOpts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseContainerOpts() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseOptional(t *testing.T) {
	type args struct {
		cmd *cobra.Command
	}
	tests := []struct {
		name           string
		args           args
		wantVolumeName string
		wantShellCmd   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVolumeName, gotShellCmd := parseOptional(tt.args.cmd)
			if gotVolumeName != tt.wantVolumeName {
				t.Errorf("parseOptional() gotVolumeName = %v, want %v", gotVolumeName, tt.wantVolumeName)
			}
			if gotShellCmd != tt.wantShellCmd {
				t.Errorf("parseOptional() gotShellCmd = %v, want %v", gotShellCmd, tt.wantShellCmd)
			}
		})
	}
}
