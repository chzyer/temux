package term

import (
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"github.com/pkg/term"
	"github.com/pkg/term/termios"
)

type Term struct {
	*term.Term
}

func RawMode(t *Term) error {
	return term.RawMode(t.Term)
}

func New(f *os.File, options ...func(*Term) error) (*Term, error) {
	t, err := term.Open(f.Name())
	if err != nil {
		return nil, err
	}
	tw := &Term{Term: t}
	if err := tw.SetOption(options...); err != nil {
		return nil, err
	}
	return tw, nil
}

func (t *Term) SetOption(options ...func(*Term) error) error {
	for _, opt := range options {
		if err := opt(t); err != nil {
			return err
		}
	}
	return nil
}

type Pty struct {
	pty *os.File
	tty *os.File
}

func NewPty() (*Pty, error) {
	pty, tty, err := termios.Pty()
	if err != nil {
		return nil, err
	}
	return &Pty{
		pty: pty,
		tty: tty,
	}, nil
}

func (p *Pty) Copy(o *os.File) (err error) {
	var attr syscall.Termios
	if err = termios.Tcgetattr(o.Fd(), &attr); err != nil {
		return err
	}
	w, h, err := GetSize(o.Fd())
	if err != nil {
		return err
	}

	if err = SetSize(p.tty.Fd(), w, h); err != nil {
		return err
	}

	err = termios.Tcsetattr(p.tty.Fd(), termios.TCSADRAIN, &attr)
	if err != nil {
		return err
	}
	return nil
}

func (p *Pty) SetCmd(c *exec.Cmd) {
	c.Stdin = p.tty
	c.Stdout = p.tty
	c.Stderr = p.tty
	if c.SysProcAttr == nil {
		c.SysProcAttr = &syscall.SysProcAttr{}
	}
	c.SysProcAttr.Setctty = true
	c.SysProcAttr.Setsid = true
}

func (p *Pty) Read(b []byte) (int, error) {
	return p.pty.Read(b)
}

func (p *Pty) Write(b []byte) (int, error) {
	return p.pty.Write(b)
}

func (p *Pty) Close() {
	p.tty.Close()
	p.pty.Close()
}

func (p *Pty) WindowChange(w, h int) error {
	return SetSize(p.pty.Fd(), w, h)
}

func GetSize(fd uintptr) (width, height int, err error) {
	var dimensions [4]uint16
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		fd, syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&dimensions)))
	if errno != 0 {
		return -1, -1, errno
	}
	return int(dimensions[1]), int(dimensions[0]), nil
}

func SetSize(fd uintptr, width, height int) (err error) {
	dimensions := [4]uint16{uint16(height), uint16(width), 0, 0}
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		fd, syscall.TIOCSWINSZ, uintptr(unsafe.Pointer(&dimensions)))
	if errno != 0 {
		return errno
	}
	return nil
}

func CopyTo(cmd *exec.Cmd, fd *os.File) (*Pty, error) {
	pty, err := NewPty()
	if err != nil {
		return nil, err
	}
	if err := pty.Copy(os.Stdin); err != nil {
		return nil, err
	}
	pty.SetCmd(cmd)
	return pty, nil
}
