package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	catLinesNumbers bool
	catLinesServer  string
)

var catLinesCmd = &cobra.Command{
	Use:   "cat-lines [archivo] [linea_inicio] [linea_fin]",
	Short: "Lee un rango específico de líneas de un archivo",
	Long: `Lee y muestra líneas específicas de un archivo remoto.
Perfecto para inspeccionar secciones de código sin leer todo el archivo.

Ejemplos:
  sshcli cat-lines /app/main.py 1 50          # Primeras 50 líneas
  sshcli cat-lines /app/main.py 100 150       # Líneas 100-150
  sshcli cat-lines -n /app/main.py 25 75      # Con números de línea
  sshcli cat-lines --server dev /app/config.py 1 20

Casos de uso para agentes:
  - Inspeccionar funciones específicas
  - Ver contexto alrededor de un error
  - Leer secciones de configuración
  - Analizar código sin cargar archivos grandes`,
	Args: cobra.ExactArgs(3),
	RunE: runCatLines,
}

func init() {
	rootCmd.AddCommand(catLinesCmd)
	catLinesCmd.Flags().BoolVarP(&catLinesNumbers, "numbers", "n", false, "Mostrar números de línea")
	catLinesCmd.Flags().StringVarP(&catLinesServer, "server", "s", "", "Servidor específico a usar")
}

func runCatLines(cmd *cobra.Command, args []string) error {
	remotePath := args[0]
	startLine, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("número de línea inicial inválido: %v", err)
	}
	endLine, err := strconv.Atoi(args[2])
	if err != nil {
		return fmt.Errorf("número de línea final inválido: %v", err)
	}

	if startLine < 1 {
		startLine = 1
	}
	if startLine > endLine {
		startLine, endLine = endLine, startLine
	}

	client, _, err := getClient(catLinesServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	data, err := client.ReadFile(remotePath)
	if err != nil {
		return fmt.Errorf("error al leer archivo: %v", err)
	}

	lines := strings.Split(string(data), "\n")
	totalLines := len(lines)

	if startLine > totalLines {
		return fmt.Errorf("línea inicial %d excede el total de líneas (%d)", startLine, totalLines)
	}

	if endLine > totalLines {
		endLine = totalLines
	}

	// Ajustar a índices 0-based
	startIdx := startLine - 1
	endIdx := endLine

	for i := startIdx; i < endIdx; i++ {
		if catLinesNumbers {
			fmt.Printf("%4d | %s\n", i+1, lines[i])
		} else {
			fmt.Println(lines[i])
		}
	}

	return nil
}
