package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var findCodeServer string

var findCodeCmd = &cobra.Command{
	Use:   "find-code [patron] [ruta]",
	Short: "Busca definiciones de código (funciones, clases) en una ruta",
	Long: `Busca patrones de código de forma recursiva.
Ideal para encontrar dónde se define una función o clase.

Ejemplos:
  sshcli find-code "def login" /app
  sshcli find-code "class User" /app/models
  sshcli find-code "export const" /app/src`,
	Args: cobra.ExactArgs(2),
	RunE: runFindCode,
}

func init() {
	rootCmd.AddCommand(findCodeCmd)
	findCodeCmd.Flags().StringVarP(&findCodeServer, "server", "s", "", "Servidor específico")
}

func runFindCode(cmd *cobra.Command, args []string) error {
	pattern := args[0]
	remotePath := cleanRemotePath(args[1])

	client, _, err := getClient(findCodeServer)
	if err != nil {
		return err
	}
	defer client.Close()

	// Busca el patrón con número de línea, ignorando binarios y carpetas .git
	grepCmd := fmt.Sprintf("grep -rnE '%s' '%s' --exclude-dir=.git --color=never", pattern, remotePath)
	output, err := client.Run(grepCmd)
	if err != nil {
		fmt.Println("No se encontraron coincidencias.")
		return nil
	}

	fmt.Print(output)
	return nil
}