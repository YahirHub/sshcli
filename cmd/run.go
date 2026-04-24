package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	runArgs      string
	runEnvVars   string
	runWorkdir   string
	runServer    string
)

var runCmd = &cobra.Command{
	Use:   "run [archivo]",
	Short: "Ejecuta un script con su intérprete apropiado",
	Long: `Ejecuta un archivo de código con el intérprete apropiado.
Detecta automáticamente el lenguaje por la extensión.

Ejemplos:
  sshcli run /app/main.py
  sshcli run /app/script.js --args "--port 3000"
  sshcli run --server prod /app/deploy.sh
  sshcli run /app/main.go --workdir /app
  sshcli run /app/test.py --env "DEBUG=1 API_KEY=xxx"

Intérpretes soportados:
  .py     -> python3
  .js     -> node
  .ts     -> ts-node / npx tsx
  .go     -> go run
  .sh     -> bash
  .rb     -> ruby
  .php    -> php
  .pl     -> perl`,
	Args: cobra.ExactArgs(1),
	RunE: runRunCmd,
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&runArgs, "args", "a", "", "Argumentos para el script")
	runCmd.Flags().StringVarP(&runEnvVars, "env", "e", "", "Variables de entorno (KEY=value KEY2=value2)")
	runCmd.Flags().StringVarP(&runWorkdir, "workdir", "w", "", "Directorio de trabajo")
	runCmd.Flags().StringVarP(&runServer, "server", "s", "", "Servidor específico a usar")
}

func runRunCmd(cmd *cobra.Command, args []string) error {
	remotePath := cleanRemotePath(args[0])

	client, _, err := getClient(runServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	ext := strings.ToLower(filepath.Ext(remotePath))
	var interpreter string

	switch ext {
	case ".py":
		interpreter = "python3"
	case ".js":
		interpreter = "node"
	case ".ts":
		interpreter = "npx tsx"
	case ".go":
		interpreter = "go run"
	case ".sh", ".bash":
		interpreter = "bash"
	case ".rb":
		interpreter = "ruby"
	case ".php":
		interpreter = "php"
	case ".pl":
		interpreter = "perl"
	default:
		return fmt.Errorf("extensión '%s' no soportada", ext)
	}

	// Construir comando
	var cmdParts []string

	// Agregar variables de entorno
	if runEnvVars != "" {
		cmdParts = append(cmdParts, runEnvVars)
	}

	// Cambiar directorio de trabajo
	if runWorkdir != "" {
		cmdParts = append(cmdParts, fmt.Sprintf("cd %s &&", runWorkdir))
	}

	// Comando principal
	cmdParts = append(cmdParts, interpreter, remotePath)

	// Agregar argumentos
	if runArgs != "" {
		cmdParts = append(cmdParts, runArgs)
	}

	fullCmd := strings.Join(cmdParts, " ")

	output, err := client.Run(fullCmd)
	if err != nil {
		fmt.Printf("Error al ejecutar %s:\n", remotePath)
		if output != "" {
			fmt.Print(output)
		}
		return err
	}

	fmt.Print(output)
	return nil
}