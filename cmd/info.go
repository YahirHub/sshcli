package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var infoServer string

var infoCmd = &cobra.Command{
	Use:   "info [ruta_remota]",
	Short: "Muestra información detallada de un archivo o directorio",
	Long: `Muestra metadatos de un archivo o directorio remoto:
tamaño, permisos, propietario, fecha de modificación.

Ejemplos:
  sshcli info /home/user/archivo.txt
  sshcli info --server prod /var/log/app.log
  sshcli info /etc/nginx`,
	Args: cobra.ExactArgs(1),
	RunE: runInfo,
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringVarP(&infoServer, "server", "s", "", "Servidor específico a usar")
}

func runInfo(cmd *cobra.Command, args []string) error {
	remotePath := args[0]

	client, _, err := getClient(infoServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	output, err := client.Run(fmt.Sprintf("stat %s 2>/dev/null || ls -la %s", remotePath, remotePath))
	if err != nil {
		return fmt.Errorf("error al obtener información: %v", err)
	}

	fmt.Print(output)
	return nil
}
