package cmd

// TODO: move this package to lib or utils folder

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	ioutils2 "github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/moby/term"
	"io"
	"log"
	"runtime"
	"sync"
)

var defaultEscapeKeys = []byte{16, 17}

type hijackedIOStreamer struct {
	streams      Streams
	inputStream  io.ReadCloser
	outputStream io.Writer
	errorStream  io.Writer

	resp types.HijackedResponse

	tty        bool
	detachKeys string
}

func (h *hijackedIOStreamer) stream(ctx context.Context) error {
	restoreInput, err := h.setupInput()
	if err != nil {
		return fmt.Errorf("unable to setup input stream: %v", err)
	}

	defer restoreInput()

	outputDone := h.beginOutputStream(restoreInput)
	inputDone, detached := h.beginInputStream(restoreInput)

	select {
	case err := <-outputDone:
		return err
	case <-inputDone:
		// input stream has closed
		if h.outputStream != nil || h.errorStream != nil {
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

func (h *hijackedIOStreamer) setupInput() (restore func(), err error) {
	if h.inputStream == nil || !h.tty {
		// no need to setup input TTY
		// the resotre func is a noop
		return func() {}, nil
	}

	// Use sync.Once so we may call restore multiple times but ensure we
	// only restore the terminal once.
	var restoreOnce sync.Once
	restore = func() {
		restoreOnce.Do(func() {
			restoreTerminal(h.streams, h.inputStream)
		})
	}

	// Wrap the input to detect detach escape sequence.
	// Use default escape keys if an invalid sequence is given.
	escapeKeys := defaultEscapeKeys
	if h.detachKeys != "" {
		customEscapeKeys, err := term.ToBytes(h.detachKeys)
		if err != nil {
			log.Printf("invalid detach escape keys, using default: %s", err)
		} else {
			escapeKeys = customEscapeKeys
		}
	}

	h.inputStream = ioutils2.NewReadCloserWrapper(term.NewEscapeProxy(h.inputStream, escapeKeys), h.inputStream.Close)

	return restore, nil
}

func (h *hijackedIOStreamer) beginOutputStream(restoreInput func()) <-chan error {
	// TODO: make output error channel and copy output from container
	if h.outputStream == nil && h.errorStream == nil {
		return nil
	}

	outputDone := make(chan error)
	go func() {
		var err error

		if h.outputStream != nil && h.tty {
			_, err = io.Copy(h.outputStream, h.resp.Reader)
			// restore terminal as soon as possible once connection ends so
			// any following print messages will be in normal type
			restoreInput()
		} else {
			_, err = stdcopy.StdCopy(h.outputStream, h.errorStream, h.resp.Reader)

			DebugPrint(fmt.Sprintf("[hijack] End of stdout"))

			if err != nil {
				DebugPrint(fmt.Sprintf("error from receiveStdout: %v", err))
			}

			outputDone <- err
		}
	}()

	return outputDone
}

func (h *hijackedIOStreamer) beginInputStream(restoreInput func()) (doneC <-chan struct{}, detachedC <-chan error) {
	inputDone := make(chan struct{})
	detached := make(chan error)

	go func() {
		if h.inputStream != nil {
			_, err := io.Copy(h.resp.Conn, h.inputStream)
			restoreInput()

			log.Printf("[hijack] End of stdin")

			if _, ok := err.(term.EscapeError); ok {
				detached <- err
				return
			}

			if err != nil {
				log.Printf("Error sendStdin: %v", err)
			}
		}

		if err := h.resp.CloseWrite(); err != nil {
			log.Printf("couldn't send EOF: %v", err)
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
