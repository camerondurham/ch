package cmd

import (
	"fmt"
	"github.com/camerondurham/ch/cmd/util"
	"github.com/docker/docker/api/types/filters"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

// runningCmd represents the running command
var runningCmd = &cobra.Command{
	Use:   "running",
	Short: "list running environments",
	Long: `List all running Docker containers created by the container-helper.
Any Docker containers not managed by the container-helper will be ignored.
To see all running containers, run: 
	docker ps`,
	Run: RunningCmd,
}

func RunningCmd(cmd *cobra.Command, args []string) {
	cli, err := util.NewCliClient()
	if err != nil {
		fmt.Printf("error: cannot create new CLI ApiClient: %v\n", err)
		os.Exit(1)
	}

	envs := cli.Containers()
	environmentNames := make([]filters.KeyValuePair, 0)
	for name, _ := range envs {
		environmentNames = append(environmentNames, filters.Arg("name", name))
	}
	f := filters.NewArgs(environmentNames...)
	list, err := cli.DockerClient().GetRunning(f, false)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintf(w, "ENVIRONMENT\tIMAGE NAME\tCREATED\n")
	for _, info := range list {
		now := time.Now().Unix()
		elapsed := now - info.Created
		var timeSinceCreated string
		if elapsed >= 60 {
			timeSinceCreated = fmt.Sprintf("%v minutes ago", elapsed/60)
		} else {
			timeSinceCreated = fmt.Sprintf("%v seconds ago", elapsed)
		}
		fmt.Fprintf(w, "%s\t%s\t%v\n", strings.TrimPrefix(info.Names[0], "/"), info.Image, timeSinceCreated)
	}
	w.Flush()
}

func init() {
	rootCmd.AddCommand(runningCmd)
}
