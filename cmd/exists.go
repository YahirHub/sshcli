package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var existsServer string

var existsCmd = &cobra.Command{
	Use:   "exists [ruta_remota]",
	Short: "Verifica si un archivo o directorio existe",
	Long: `Verifica si un archivo o directorio existe en el servidor remoto.
Retorna código de salida 0 si existe, 1 si no existe.
Ideal para agentes que necesitan verificar antes de crear/modificar.

Ejemplos:
  sshcli exists /home/user/archivo.txt
  sshcli exists --server prod /var/www/app
  sshcli exists /etc/nginx/nginx.conf && echo "Existe"`,
	Args: cobra.ExactArgs(1),
	RunE: runExists,
}

func init() {
	rootCmd.AddCommand(existsCmd)
	existsCmd.Flags().StringVarP(&existsServer, "server", "s", "", "Servidor específico a usar")
}

func runExists(cmd *cobra.Command, args []string) error {
	remotePath := args[0]

	client, _, err := getClient(existsServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	_, err = client.Run(fmt.Sprintf("test -e %s", remotePath))
	if err != nil {
		fmt.Println("NO")
		return fmt.Errorf("no existe: %s", remotePath)
	}

	fmt.Println("SI")
	return nil
}
