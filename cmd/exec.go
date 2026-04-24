package cmd

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"

	"sshcli/internal/config"

	"github.com/spf13/cobra"
)

var (
	execServer string
	execTTY    bool
	execNoTTY  bool
)

var execCmd = &cobra.Command{
	Use:   "exec [comando]",
	Short: "Ejecuta un comando en el servidor remoto",
	Long: `Ejecuta un comando en el servidor remoto configurado.`,
	Args: cobra.MinimumNArgs(1),
	RunE: runExec,
}

func init() {
	rootCmd.AddCommand(execCmd)
	execCmd.Flags().StringVarP(&execServer, "server", "s", "", "Servidor específico a usar")
	execCmd.Flags().BoolVarP(&execTTY, "tty", "t", false, "Habilitar modo interactivo (PTY)")
	execCmd.Flags().BoolVar(&execNoTTY, "no-tty", false, "Forzar modo normal (ignora config)")
}

func runExec(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("error de configuración: %v", err)
	}

	var server *config.Server
	if execServer != "" {
		server, err = cfg.GetServer(execServer)
	} else {
		server, err = cfg.GetActiveServer()
	}
	if err != nil {
		return err
	}

	// Simplificado: No intentamos limpiar comandos complejos automáticamente
	// para evitar añadir slashes innecesarios al principio.
	command := strings.Join(args, " ")

	sshConfig := &ssh.ClientConfig{
		User: server.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(server.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%d", server.Host, server.Port)
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return fmt.Errorf("error de conexión: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("error de sesión: %v", err)
	}
	defer session.Close()

	useTTY := (execTTY || cfg.DefaultTTY) && !execNoTTY

	if useTTY {
		return runInteractiveExec(session, command)
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	return session.Run(command)
}

func runInteractiveExec(session *ssh.Session, command string) error {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		session.Stdout = os.Stdout
		session.Stderr = os.Stderr
		return session.Run(command)
	}

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer term.Restore(fd, oldState)

	width, height, err := term.GetSize(fd)
	if err != nil {
		width, height = 80, 24
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm-256color", height, width, modes); err != nil {
		return err
	}

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	return session.Run(command)
}