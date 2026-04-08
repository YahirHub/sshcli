//go:build !windows

package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func handleWindowResize(session *ssh.Session, fd int) {
	sigwinch := make(chan os.Signal, 1)
	signal.Notify(sigwinch, syscall.SIGWINCH)
	for range sigwinch {
		if w, h, err := term.GetSize(fd); err == nil {
			session.WindowChange(h, w)
		}
	}
}
