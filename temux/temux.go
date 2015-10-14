package main

import (
	"io"
	"os"
	"os/exec"

	"github.com/chzyer/temux/temux/term"
	"github.com/pkg/term/termios"
)

func main() {
	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	pty, tty, err := termios.Pty()
	if err != nil {
		panic(err)
	}
	cmd.Stdin = tty
	cmd.Stdout = tty
	cmd.Stderr = tty
	if err := cmd.Start(); err != nil {
		return
	}
	tty.Close()

	t, err := term.New(os.Stdin, term.RawMode)
	if err != nil {
		panic(err)
	}
	defer t.Restore()

	go io.Copy(os.Stdout, pty)
	go io.Copy(pty, os.Stdin)
	cmd.Wait()
}
