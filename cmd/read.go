package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var readServer string

var readCmd = &cobra.Command{
	Use:   "read [ruta_remota]",
	Short: "Lee el contenido de un archivo remoto",
	Long: `Lee y muestra el contenido de un archivo remoto.
Ideal para agentes de IA que necesitan verificar el contenido
de archivos en el servidor.

Ejemplos:
  sshcli read /home/user/archivo.txt
  sshcli read --server produccion /etc/nginx/nginx.conf
  sshcli read /var/log/app.log`,
	Args: cobra.ExactArgs(1),
	RunE: runRead,
}

func init() {
	rootCmd.AddCommand(readCmd)
	readCmd.Flags().StringVarP(&readServer, "server", "s", "", "Servidor específico a usar")
}

func runRead(cmd *cobra.Command, args []string) error {
	remotePath := args[0]

	client, _, err := getClient(readServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	data, err := client.ReadFile(remotePath)
	if err != nil {
		return fmt.Errorf("error al leer archivo: %v", err)
	}

	fmt.Print(string(data))
	return nil
}
