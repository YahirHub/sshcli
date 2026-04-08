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
	Long: `Ejecuta un comando en el servidor remoto configurado.
El comando puede ser cualquier instrucción válida de shell.

MODO NORMAL (sin -t):
  sshcli exec "ls -la"
  sshcli exec "cat /etc/hostname"
  sshcli exec "echo hello"

MODO INTERACTIVO (-t):
  sshcli exec -t htop              # Programas de pantalla completa
  sshcli exec -t "apt install x"   # Confirmaciones Y/N
  sshcli exec -t "vim archivo"     # Editores
  sshcli exec -t bash              # Shell interactivo

CONFIGURACIÓN:
  Puedes habilitar -t por defecto con:
    sshcli config set tty true

  Y desactivarlo temporalmente con --no-tty:
    sshcli exec --no-tty "ls -la"

NOTAS:
  - Usa -t para programas que requieren entrada de teclado
  - Sin -t para scripts y comandos simples (más rápido)
  - Si el terminal se bugea, escribe 'reset' y presiona Enter`,
	Args: cobra.MinimumNArgs(1),
	RunE: runExec,
}

func init() {
	rootCmd.AddCommand(execCmd)
	execCmd.Flags().StringVarP(&execServer, "server", "s", "", "Servidor específico a usar")
	execCmd.Flags().BoolVarP(&execTTY, "tty", "t", false, "Modo interactivo con pseudo-terminal")
	execCmd.Flags().BoolVar(&execNoTTY, "no-tty", false, "Forzar modo no interactivo")
}

func runExec(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("error: %v. Usa 'sshcli server add' para configurar", err)
	}

	var server *config.Server
	if execServer != "" {
		server, err = cfg.GetServer(execServer)
	} else {
		server, err = cfg.GetActiveServer()
	}
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	// Determinar si usar TTY
	useTTY := execTTY || cfg.DefaultTTY
	if execNoTTY {
		useTTY = false
	}

	// Configurar cliente SSH
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
		return fmt.Errorf("error al crear sesión: %v", err)
	}
	defer session.Close()

	command := strings.Join(args, " ")

	if useTTY {
		return runInteractiveExec(session, command)
	}

	// Modo no interactivo
	output, err := session.CombinedOutput(command)
	if err != nil {
		if len(output) > 0 {
			fmt.Print(string(output))
		}
		return fmt.Errorf("error al ejecutar comando: %v", err)
	}
	fmt.Print(string(output))
	return nil
}

func runInteractiveExec(session *ssh.Session, command string) error {
	fd := int(os.Stdin.Fd())
	var oldState *term.State
	var err error

	// Solo configurar raw mode si es un terminal real
	if term.IsTerminal(fd) {
		oldState, err = term.MakeRaw(fd)
		if err != nil {
			// Si falla, continuar sin raw mode
			oldState = nil
		}
	}

	// SIEMPRE restaurar el terminal al salir
	defer func() {
		if oldState != nil {
			term.Restore(fd, oldState)
		}
	}()

	// Obtener tamaño de terminal
	width, height := 80, 24
	if term.IsTerminal(fd) {
		if w, h, err := term.GetSize(fd); err == nil {
			width, height = w, h
		}
	}

	// Configurar PTY remoto
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm-256color", height, width, modes); err != nil {
		return fmt.Errorf("error al solicitar PTY: %v", err)
	}

	// Conectar stdin, stdout, stderr
	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	// Ejecutar comando
	err = session.Run(command)
	
	// Restaurar terminal ANTES de retornar error
	if oldState != nil {
		term.Restore(fd, oldState)
		oldState = nil // Evitar doble restore en defer
	}

	if err != nil {
		if _, ok := err.(*ssh.ExitError); ok {
			return nil
		}
		return err
	}

	return nil
}
