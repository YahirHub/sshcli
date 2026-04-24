package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var insertLineServer string

var insertLineCmd = &cobra.Command{
	Use:   "insert-line [archivo] [numero_linea] [contenido]",
	Short: "Inserta una línea en una posición específica del archivo",
	Long: `Inserta contenido en una línea específica de un archivo remoto.
La línea existente y las siguientes se desplazan hacia abajo.
Usa línea 0 para insertar al inicio del archivo.

Ejemplos:
  sshcli insert-line /app/main.py 1 "import os"
  sshcli insert-line /app/main.py 10 "    # TODO: refactorizar"
  sshcli insert-line --server dev /app/config.py 5 "DEBUG = True"
  
Casos de uso para agentes:
  - Agregar imports al inicio
  - Insertar código en función existente
  - Agregar comentarios
  - Insertar configuraciones`,
	Args: cobra.ExactArgs(3),
	RunE: runInsertLine,
}

func init() {
	rootCmd.AddCommand(insertLineCmd)
	insertLineCmd.Flags().StringVarP(&insertLineServer, "server", "s", "", "Servidor específico a usar")
}

func runInsertLine(cmd *cobra.Command, args []string) error {
	remotePath := cleanRemotePath(args[0])
	lineNum, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("número de línea inválido: %v", err)
	}
	content := args[2]

	client, _, err := getClient(insertLineServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	// Leer archivo
	data, err := client.ReadFile(remotePath)
	if err != nil {
		return fmt.Errorf("error al leer archivo: %v", err)
	}

	lines := strings.Split(string(data), "\n")
	
	// Validar número de línea
	if lineNum < 0 {
		lineNum = 0
	}
	if lineNum > len(lines) {
		lineNum = len(lines)
	}

	// Insertar línea
	newLines := make([]string, 0, len(lines)+1)
	newLines = append(newLines, lines[:lineNum]...)
	newLines = append(newLines, content)
	newLines = append(newLines, lines[lineNum:]...)

	newContent := strings.Join(newLines, "\n")

	if err := client.WriteFile(remotePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("error al escribir archivo: %v", err)
	}

	fmt.Printf("Línea insertada en posición %d: %s\n", lineNum+1, remotePath)
	return nil
}