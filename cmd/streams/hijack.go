package streams

// TODO: move this package to lib or utils folder

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	ioutils2 "github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/moby/term"
	"io"
	"os"
	"runtime"
	"sync"
)

var defaultEscapeKeys = []byte{16, 17}

// Streams is an interface which exposes the standard input and output streams
type Streams interface {
	In() *In
	Out() *Out
	Err() io.Writer
}

type HijackedIOStreamer struct {
	Streams      Streams
	InputStream  io.ReadCloser
	OutputStream io.Writer
	ErrorStream  io.Writer

	Resp types.HijackedResponse

	Tty        bool
	detachKeys string
}

func (h *HijackedIOStreamer) Stream(ctx context.Context) error {
	restoreInput, err := h.setupInput()
	if err != nil {
		return fmt.Errorf("unable to setup input Stream: %v", err)
	}

	defer restoreInput()

	outputDone := h.beginOutputStream(restoreInput)
	inputDone, detached := h.beginInputStream(restoreInput)

	select {
	case err := <-outputDone:
		return err
	case <-inputDone:
		// input Stream has closed
		if h.OutputStream != nil || h.ErrorStream != nil {
			// wait for output to complete streaming.
			select {
			case err := <-outputDone:
				return err
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return nil
	case err := <-detached:
		// Got detached key sequence
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (h *HijackedIOStreamer) setupInput() (restore func(), err error) {
	if h.InputStream == nil || !h.Tty {
		// no need to setup input TTY
		// the resotre func is a noop
		return func() {}, nil
	}

	if err := setRawTerminal(h.Streams); err != nil {
		return nil, fmt.Errorf("unable to set IO Streams as raw terminal: %v", err)
	}

	// Use sync.Once so we may call restore multiple times but ensure we
	// only restore the terminal once.
	var restoreOnce sync.Once
	restore = func() {
		restoreOnce.Do(func() {
			// this function executes when user exits shell by typing `exit` or `^D`
			if err = restoreTerminal(h.Streams, h.InputStream); err != nil {
				debugPrint("error restoring terminal")
			} else {
				// Goroutine does not exit if this doesn't explicitly terminate the program.
				// Pending investigation, this should terminate the program once the user exits terminal.
				// Terminating here should be safe as the open file descriptor (stdin) is restored, resources
				// have already been restored at this point.
				os.Exit(0)
			}
		})
	}

	// Wrap the input to detect detach escape sequence.
	// Use default escape keys if an invalid sequence is given.
	escapeKeys := defaultEscapeKeys
	if h.detachKeys != "" {
		customEscapeKeys, err := term.ToBytes(h.detachKeys)
		if err != nil {
			debugPrint(fmt.Sprintf("invalid detach escape keys, using default: %s", err))
		} else {
			escapeKeys = customEscapeKeys
		}
	}

	h.InputStream = ioutils2.NewReadCloserWrapper(term.NewEscapeProxy(h.InputStream, escapeKeys), h.InputStream.Close)

	return restore, nil
}

func (h *HijackedIOStreamer) beginOutputStream(restoreInput func()) <-chan error {
	if h.OutputStream == nil && h.ErrorStream == nil {
		return nil
	}

	outputDone := make(chan error)
	go func() {
		var err error

		if h.OutputStream != nil && h.Tty {
			_, err = io.Copy(h.OutputStream, h.Resp.Reader)
			// restore terminal as soon as possible once connection ends so
			// any following print messages will be in normal type
			restoreInput()
		} else {
			_, err = stdcopy.StdCopy(h.OutputStream, h.ErrorStream, h.Resp.Reader)

			debugPrint(fmt.Sprintf("[hijack] End of stdout"))

			if err != nil {
				debugPrint(fmt.Sprintf("error from receiveStdout: %v", err))
			}

			outputDone <- err
		}
	}()

	return outputDone
}

func (h *HijackedIOStreamer) beginInputStream(restoreInput func()) (doneC <-chan struct{}, detachedC <-chan error) {
	inputDone := make(chan struct{})
	detached := make(chan error)

	go func() {
		if h.InputStream != nil {
			_, err := io.Copy(h.Resp.Conn, h.InputStream)
			restoreInput()

			debugPrint("\n[hijack] End of stdin\n")

			if _, ok := err.(term.EscapeError); ok {
				detached <- err
				return
			}

			if err != nil {
				debugPrint(fmt.Sprintf("Error sendStdin: %v", err))
			}
		}

		if err := h.Resp.CloseWrite(); err != nil {
			debugPrint(fmt.Sprintf("couldn't send EOF: %v", err))
		}
		close(inputDone)
	}()
	return inputDone, detached
}

func setRawTerminal(streams Streams) error {
	if err := streams.In().SetRawTerminal(); err != nil {
		return err
	}
	return streams.Out().SetRawTerminal()
}

func restoreTerminal(streams Streams, in io.Closer) error {
	streams.In().RestoreTerminal()
	streams.Out().RestoreTerminal()
	// WARNING: DO NOT REMOVE THE OS CHECKS !!!
	// For some reason this Close call blocks on darwin..
	// As the client exits right after, simply discard the close
	// until we find a better solution.
	//
	// This can also cause the client on Windows to get stuck in Win32 CloseHandle()
	// in some cases. See https://github.com/docker/docker/issues/28267#issuecomment-288237442
	// Tracked internally at Microsoft by VSO #11352156. In the
	// Windows case, you hit this if you are using the native/v2 console,
	// not the "legacy" console, and you start the client in a new window. eg
	// `start docker run --rm -it microsoft/nanoserver cmd /s /c echo foobar`
	// will hang. Remove start, and it won't repro.
	if in != nil && runtime.GOOS != "darwin" && runtime.GOOS != "windows" {
		return in.Close()
	}
	return nil
}

// DebugPrint if DEBUG environment variable is set
func debugPrint(msg string) {
	if _, ok := os.LookupEnv("DEBUG"); ok {
		fmt.Printf("[hijack] %v\n", msg)
	}
}
