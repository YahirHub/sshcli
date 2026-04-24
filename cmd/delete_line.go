package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var deleteLineServer string

var deleteLineCmd = &cobra.Command{
	Use:   "delete-line [archivo] [linea_inicio] [linea_fin]",
	Short: "Elimina líneas específicas de un archivo",
	Long: `Elimina una o más líneas de un archivo remoto.
Especifica una sola línea o un rango (inicio-fin).

Ejemplos:
  sshcli delete-line /app/main.py 15 15        # Elimina línea 15
  sshcli delete-line /app/main.py 10 20        # Elimina líneas 10-20
  sshcli delete-line --server prod /app/config.py 5 5
  
Casos de uso para agentes:
  - Eliminar código obsoleto
  - Quitar imports no usados
  - Borrar comentarios
  - Limpiar líneas en blanco`,
	Args: cobra.RangeArgs(2, 3),
	RunE: runDeleteLine,
}

func init() {
	rootCmd.AddCommand(deleteLineCmd)
	deleteLineCmd.Flags().StringVarP(&deleteLineServer, "server", "s", "", "Servidor específico a usar")
}

func runDeleteLine(cmd *cobra.Command, args []string) error {
	remotePath := cleanRemotePath(args[0])
	startLine, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("número de línea inicial inválido: %v", err)
	}

	endLine := startLine
	if len(args) == 3 {
		endLine, err = strconv.Atoi(args[2])
		if err != nil {
			return fmt.Errorf("número de línea final inválido: %v", err)
		}
	}

	if startLine > endLine {
		startLine, endLine = endLine, startLine
	}

	client, _, err := getClient(deleteLineServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	data, err := client.ReadFile(remotePath)
	if err != nil {
		return fmt.Errorf("error al leer archivo: %v", err)
	}

	lines := strings.Split(string(data), "\n")
	
	// Ajustar a índices 0-based
	startIdx := startLine - 1
	endIdx := endLine

	if startIdx < 0 {
		startIdx = 0
	}
	if endIdx > len(lines) {
		endIdx = len(lines)
	}

	// Eliminar líneas
	newLines := make([]string, 0, len(lines)-(endIdx-startIdx))
	newLines = append(newLines, lines[:startIdx]...)
	newLines = append(newLines, lines[endIdx:]...)

	newContent := strings.Join(newLines, "\n")

	if err := client.WriteFile(remotePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("error al escribir archivo: %v", err)
	}

	deleted := endIdx - startIdx
	if deleted == 1 {
		fmt.Printf("Eliminada línea %d: %s\n", startLine, remotePath)
	} else {
		fmt.Printf("Eliminadas %d líneas (%d-%d): %s\n", deleted, startLine, endLine, remotePath)
	}
	return nil
}