package cmd

import (
	"fmt"
	"github.com/camerondurham/ch/cmd/util"
	"github.com/spf13/cobra"
	"runtime"
)

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Short:   "instructions to upgrade the ch cli",
	Version: rootCmd.Version,
	Run:     UpgradeCmd,
}

const UnixUpgradeTerminalCommand = "bash <(curl -s https://raw.githubusercontent.com/camerondurham/ch/main/scripts/install-ch.sh)"
const WindowsUpgradeTerminalCommand = "Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/camerondurham/ch/main/scripts/install-ch.ps1'))"

func UpgradeCmd(cmd *cobra.Command, args []string) {

	var upgradeCommand string
	var os string

	switch runtime.GOOS {
	case "darwin":
		os = "macOS"
		fallthrough
	case "linux":
		// set upgrade command
		os = "Linux"
		upgradeCommand = UnixUpgradeTerminalCommand
	case "windows":
		// set upgrade command
		os = "Windows"
		upgradeCommand = WindowsUpgradeTerminalCommand
	}

	fmt.Printf("You appear to be running a %s operating system."+
		"\nPlease run the following upgrade command:"+
		"\n\n    %s"+
		"\n\nFor more help, see the repository README: %s\n", os, upgradeCommand, util.RepositoryUrl)

}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
