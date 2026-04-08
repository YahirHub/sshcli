//go:build windows

package cmd

import "golang.org/x/crypto/ssh"

func handleWindowResize(session *ssh.Session, fd int) {
	// No soportado en Windows
}
