/*
Copyright © 2020 Cameron Durham <cameron.r.durham@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"github.com/camerondurham/ch/cmd/util"
	"github.com/spf13/cobra"
	"os"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [ENVIRONMENT_NAME]",
	Short: "list configuration for existing environments",
	Run:   ListCmd,
}

func ListCmd(cmd *cobra.Command, args []string) {

	envName := ""
	if len(args) > 0 {
		envName = args[0]
	}

	configEnvs, err := util.GetEnvs()

	if err != nil {
		if err == util.ErrDoesNotExist {
			fmt.Printf("no environments created yet")
			os.Exit(0)
		}
		fmt.Printf("unable to decode config file: %v\nplease check formatting of config and delete if needed", err)
		os.Exit(1)
	}

	if envName != "" {
		if v, ok := configEnvs[envName]; ok {
			util.PrintConfig(envName, v)
		} else {
			fmt.Printf("no environment found")
		}
	} else {
		for k, v := range configEnvs {
			util.PrintConfig(k, v)
		}
	}

}
func init() {
	rootCmd.AddCommand(listCmd)
}
