package main

import (
	"io"
	"os"
	"os/exec"

	"os/signal"
	"syscall"

	"github.com/chzyer/flagx"
	"github.com/chzyer/temux/temux/term"
)

type Config struct {
	Shell string
}

func NewConfig() *Config {
	var cfg Config
	flagx.Parse(&cfg)
	return &cfg
}

func main() {
	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	pty, err := term.CopyTo(cmd, os.Stdin)
	if err != nil {
		panic(err)
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	t, err := term.New(os.Stdin, term.RawMode)
	if err != nil {
		panic(err)
	}
	defer t.Restore()

	go io.Copy(os.Stdout, pty)
	go io.Copy(pty, os.Stdin)

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGWINCH)
		for {
			select {
			case <-ch:
			}
			if w, h, err := term.GetSize(os.Stdin.Fd()); err == nil {
				pty.WindowChange(w, h)
				cmd.Process.Signal(syscall.SIGWINCH)
			}
		}

	}()
	cmd.Wait()
}
