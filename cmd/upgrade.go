package cmd

import (
	"fmt"
	"github.com/camerondurham/ch/cmd/util"
	"github.com/camerondurham/ch/version"
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

const (
	UnixUpgradeTerminalCommand = "bash <(curl -s https://raw.githubusercontent.com/camerondurham/ch/main/scripts/install-ch.sh)"
	WindowsUpgradeTerminalCommand = "Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/camerondurham/ch/main/scripts/install-ch.ps1'))"
)

type upgradeTuple struct {
	OperatingSystem string
	Command string
}

type versionInfo struct {
	NewVersionAvailable bool
	VersionInfoMessage string
}

func getUpgradeCommandMap() map[string] upgradeTuple {
	return map[string] upgradeTuple{
		"darwin" : {
			OperatingSystem: "macOS",
			Command:         UnixUpgradeTerminalCommand,
		},
		"linux" : {
			OperatingSystem: "Linux",
			Command:         UnixUpgradeTerminalCommand,
		},
		"windows": {
			OperatingSystem: "Windows",
			Command: WindowsUpgradeTerminalCommand,
		},
	}
}

func newVersionAvailable() versionInfo {
	latestVersion, err := util.GetLatestVersion(util.GetRequest, util.GetGithubAPILatestReleaseURL(util.RepositoryName))
	if err != nil {
		util.DebugPrint(fmt.Sprintf("ignoring version check since error occured when retrieving latest version: %v\n", err))
		return versionInfo{
			NewVersionAvailable: false,
			VersionInfoMessage: "could not check latest version from Github",
		}
	} else if version.PkgVersion != "" && latestVersion != version.PkgVersion {
		return versionInfo{
			NewVersionAvailable: true,
			VersionInfoMessage: fmt.Sprintf( "current version %s, latest version %s", version.PkgVersion, latestVersion),
		}

	} else {
		return versionInfo{
			NewVersionAvailable: false,
			VersionInfoMessage: fmt.Sprintf("up to date, latest version is: %s", latestVersion),
		}
	}
}

func UpgradeCmd(cmd *cobra.Command, args []string) {

	latestVersionInfo := newVersionAvailable()
	fmt.Printf("\n\t%s\n", latestVersionInfo.VersionInfoMessage)

	if true {
		upgradeInfo := getUpgradeCommandMap()[runtime.GOOS]
		fmt.Printf( "\n\tPlease run the following upgrade command (detected %s operating system):"+
			"\n\n\t%s"+
			"\n\n\tFor more help, see the repository README: %s\n", upgradeInfo.OperatingSystem, upgradeInfo.Command, util.RepositoryUrl)
	}
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
