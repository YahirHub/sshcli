package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	headLines  int
	headServer string
)

var headCmd = &cobra.Command{
	Use:   "head [archivo]",
	Short: "Muestra las primeras líneas de un archivo",
	Long: `Muestra las primeras líneas de un archivo remoto.
Útil para inspeccionar el inicio de archivos de código o configuración.

Ejemplos:
  sshcli head /app/main.py
  sshcli head -n 50 /app/main.py
  sshcli head --server dev /etc/nginx/nginx.conf

Casos de uso para agentes:
  - Ver imports y dependencias
  - Inspeccionar cabeceras de archivos
  - Revisar configuración inicial
  - Verificar shebang y metadatos`,
	Args: cobra.ExactArgs(1),
	RunE: runHead,
}

func init() {
	rootCmd.AddCommand(headCmd)
	headCmd.Flags().IntVarP(&headLines, "lines", "n", 20, "Número de líneas a mostrar")
	headCmd.Flags().StringVarP(&headServer, "server", "s", "", "Servidor específico a usar")
}

func runHead(cmd *cobra.Command, args []string) error {
	remotePath := args[0]

	client, _, err := getClient(headServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	headCommand := fmt.Sprintf("head -n %s %s", strconv.Itoa(headLines), remotePath)

	output, err := client.Run(headCommand)
	if err != nil {
		return fmt.Errorf("error al leer archivo: %v", err)
	}

	fmt.Print(output)
	return nil
}
