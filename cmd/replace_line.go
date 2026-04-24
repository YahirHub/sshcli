package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var replaceLineServer string

var replaceLineCmd = &cobra.Command{
	Use:   "replace-line [archivo] [numero_linea] [nuevo_contenido]",
	Short: "Reemplaza una línea completa en un archivo",
	Long: `Reemplaza el contenido de una línea específica en un archivo remoto.
La línea existente se sustituye completamente por el nuevo contenido.

Ejemplos:
  sshcli replace-line /app/config.py 5 "DEBUG = False"
  sshcli replace-line /app/main.py 1 "#!/usr/bin/env python3"
  sshcli replace-line --server prod /etc/nginx/nginx.conf 10 "    server_name prod.example.com;"
  
Casos de uso para agentes:
  - Cambiar valores de configuración
  - Corregir líneas con errores
  - Actualizar shebang
  - Modificar declaraciones`,
	Args: cobra.ExactArgs(3),
	RunE: runReplaceLine,
}

func init() {
	rootCmd.AddCommand(replaceLineCmd)
	replaceLineCmd.Flags().StringVarP(&replaceLineServer, "server", "s", "", "Servidor específico a usar")
}

func runReplaceLine(cmd *cobra.Command, args []string) error {
	remotePath := cleanRemotePath(args[0])
	lineNum, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("número de línea inválido: %v", err)
	}
	newContent := args[2]

	if lineNum < 1 {
		return fmt.Errorf("el número de línea debe ser >= 1")
	}

	client, _, err := getClient(replaceLineServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	data, err := client.ReadFile(remotePath)
	if err != nil {
		return fmt.Errorf("error al leer archivo: %v", err)
	}

	lines := strings.Split(string(data), "\n")
	lineIdx := lineNum - 1

	if lineIdx >= len(lines) {
		return fmt.Errorf("línea %d no existe (archivo tiene %d líneas)", lineNum, len(lines))
	}

	oldContent := lines[lineIdx]
	lines[lineIdx] = newContent

	result := strings.Join(lines, "\n")

	if err := client.WriteFile(remotePath, []byte(result), 0644); err != nil {
		return fmt.Errorf("error al escribir archivo: %v", err)
	}

	fmt.Printf("Línea %d reemplazada en %s\n", lineNum, remotePath)
	fmt.Printf("  Antes: %s\n", oldContent)
	fmt.Printf("  Ahora: %s\n", newContent)
	return nil
}