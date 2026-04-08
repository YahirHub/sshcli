package cmd

import (
	"fmt"
	"os"
	"runtime"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"

	"sshcli/internal/config"

	"github.com/spf13/cobra"
)

var connectServer string

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Abre una sesión SSH interactiva",
	Long: `Abre una sesión de terminal SSH interactiva con el servidor.
Permite ejecutar comandos manualmente como en una terminal normal.
Usa Ctrl+D o 'exit' para salir.

Ejemplos:
  sshcli connect                      # Conectar al servidor activo
  sshcli connect --server produccion  # Conectar a servidor específico

Casos de uso para agentes:
  - Debugging interactivo
  - Tareas que requieren input manual
  - Exploración del servidor`,
	RunE: runConnect,
}

func init() {
	rootCmd.AddCommand(connectCmd)
	connectCmd.Flags().StringVarP(&connectServer, "server", "s", "", "Servidor específico a usar")
}

func runConnect(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("configuración no encontrada: %v", err)
	}

	var server *config.Server
	if connectServer != "" {
		server, err = cfg.GetServer(connectServer)
	} else {
		server, err = cfg.GetActiveServer()
	}
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	fmt.Printf("Conectando a %s@%s:%d...\n", server.User, server.Host, server.Port)

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

	// Crear sesión
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("error al crear sesión: %v", err)
	}
	defer session.Close()

	// Configurar terminal
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("error al configurar terminal: %v", err)
	}
	defer term.Restore(fd, oldState)

	// Obtener tamaño de terminal
	width, height, err := term.GetSize(fd)
	if err != nil {
		width, height = 80, 24
	}

	// Configurar PTY
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm-256color", height, width, modes); err != nil {
		return fmt.Errorf("error al solicitar PTY: %v", err)
	}

	// Conectar stdin/stdout/stderr
	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	// Manejar cambios de tamaño de ventana (solo en Unix)
	if runtime.GOOS != "windows" {
		go handleWindowResize(session, fd)
	}

	// Iniciar shell
	if err := session.Shell(); err != nil {
		return fmt.Errorf("error al iniciar shell: %v", err)
	}

	fmt.Printf("Conectado. Usa 'exit' o Ctrl+D para salir.\n\n")

	// Esperar a que termine la sesión
	session.Wait()

	fmt.Println("\nConexión cerrada.")
	return nil
}
