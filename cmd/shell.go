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
	"bufio"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
	"log"
	"os"
)

const (
	autostartFlag      = "force-start"
	autostartFlagShort = "f"
)

// shellCmd represents the shell command
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		envName := args[0]
		autostart, _ := cmd.Flags().GetBool(autostartFlagShort)

		// TODO: create helper function for envName existing
		envs := GetEnvsOrDie()
		if containerOpts, ok := envs[envName]; ok {
			ctx, cli := DockerClientInitOrDie()
			running, err := GetRunning()
			containerID, ok := running[envName]
			if !autostart && (err == ErrDoesNotExist || !ok) {
				fmt.Printf(getNotRunningMsg(envName))
				os.Exit(1)
			} else if err == ErrDoesNotExist || !ok {
				// TODO: start container
				fmt.Print("error: force starting not implemented yet")
				os.Exit(1)
			}

			DebugPrint(fmt.Sprintf("starting container: %v", containerID))
			execID, reader, writer, err := CreateExecInteractive(ctx, cli, containerID, types.ExecConfig{
				Cmd:          []string{containerOpts.Shell},
				Tty:          true,
				AttachStdin:  true,
				AttachStdout: true,
				AttachStderr: true,
			})

			log.Printf("created execID: %v", execID)
			log.Printf("reader: %v, writer: %v", reader, writer)

			go func() {
				buf := make([]byte, 1024)
				for {
					n, err := reader.Read(buf)
					if err != nil {
						fmt.Printf("failed to read: %v", err)
						break
					} else {
						// TODO: debug print
						log.Printf("read %d bytes: %v", n, buf[:n])
						fmt.Printf("%s", buf[:n])
					}
				}
			}()

			stdReader := bufio.NewReader(os.Stdin)
			for {
				text, err := stdReader.ReadString('\n')
				if err != nil {
					fmt.Printf("failed to read: %v", err)
				} else {
					_, err := writer.Write([]byte(text))
					if err != nil {
						fmt.Printf("could not write bytes: %v", err)
					}
				}
			}

		} else {
			fmt.Printf("no such environment: %v", envName)
		}
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
	shellCmd.Flags().BoolP(autostartFlag, autostartFlagShort, false, "autostart the environment if not running")
}

func getNotRunningMsg(envName string) string {
	return fmt.Sprintf(`%v is not running, please run: 
	ch create %v

or start container automatically with:
	ch shell %v %v\n\n`, envName, envName, envName, autostartFlag)
}
