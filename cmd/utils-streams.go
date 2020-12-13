package cmd

import (
	"errors"
	"github.com/moby/term"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"runtime"
)

/////////////////////////////////////////////////////////////////////////////////////////////////////////////

// commonStream is an input stream used by the DockerCli to read user input
type commonStream struct {
	fd         uintptr
	isTerminal bool
	state      *term.State
}

// FD returns the file descriptor number for this stream
func (s *commonStream) FD() uintptr {
	return s.fd
}

// IsTerminal returns true if this stream is connected to a terminal
func (s *commonStream) IsTerminal() bool {
	return s.isTerminal
}

// RestoreTerminal restores normal mode to the terminal
func (s *commonStream) RestoreTerminal() {
	if s.state != nil {
		term.RestoreTerminal(s.fd, s.state)
	}
}

// SetIsTerminal sets the boolean used for isTerminal
func (s *commonStream) SetIsTerminal(isTerminal bool) {
	s.isTerminal = isTerminal
}

//////////////////////////////////////////////////////////////////////////////////////

// Out is an output stream used by the DockerCli to write normal program
// output.
type Out struct {
	commonStream
	out io.Writer
}

func (o *Out) Write(p []byte) (int, error) {
	return o.out.Write(p)
}

// SetRawTerminal sets raw mode on the input terminal
func (o *Out) SetRawTerminal() (err error) {
	if os.Getenv("NORAW") != "" || !o.commonStream.isTerminal {
		return nil
	}
	o.commonStream.state, err = term.SetRawTerminalOutput(o.commonStream.fd)
	return err
}

// GetTtySize returns the height and width in characters of the tty
func (o *Out) GetTtySize() (uint, uint) {
	if !o.isTerminal {
		return 0, 0
	}
	ws, err := term.GetWinsize(o.fd)
	if err != nil {
		logrus.Debugf("Error getting size: %s", err)
		if ws == nil {
			return 0, 0
		}
	}
	return uint(ws.Height), uint(ws.Width)
}

// NewOut returns a new Out object from a Writer
func NewOut(out io.Writer) *Out {
	fd, isTerminal := term.GetFdInfo(out)
	return &Out{commonStream: commonStream{fd: fd, isTerminal: isTerminal}, out: out}
}

////////////////////////////////////////////////////////////////

// In is an input stream used by the DockerCli to read user input
type In struct {
	commonStream
	in io.ReadCloser
}

func (i *In) Read(p []byte) (int, error) {
	return i.in.Read(p)
}

// Close implements the Closer interface
func (i *In) Close() error {
	return i.in.Close()
}

// SetRawTerminal sets raw mode on the input terminal
func (i *In) SetRawTerminal() (err error) {
	if os.Getenv("NORAW") != "" || !i.commonStream.isTerminal {
		return nil
	}
	i.commonStream.state, err = term.SetRawTerminal(i.commonStream.fd)
	return err
}

// CheckTty checks if we are trying to attach to a container tty
// from a non-tty client input stream, and if so, returns an error.
func (i *In) CheckTty(attachStdin, ttyMode bool) error {
	// In order to attach to a container tty, input stream for the client must
	// be a tty itself: redirecting or piping the client standard input is
	// incompatible with `docker run -t`, `docker exec -t` or `docker attach`.
	if ttyMode && attachStdin && !i.isTerminal {
		eText := "the input device is not a TTY"
		if runtime.GOOS == "windows" {
			return errors.New(eText + ".  If you are using mintty, try prefixing the command with 'winpty'")
		}
		return errors.New(eText)
	}
	return nil
}

// NewIn returns a new In object from a ReadCloser
func NewIn(in io.ReadCloser) *In {
	fd, isTerminal := term.GetFdInfo(in)
	return &In{commonStream: commonStream{fd: fd, isTerminal: isTerminal}, in: in}
}
