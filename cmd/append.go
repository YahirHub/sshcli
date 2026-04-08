package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var appendServer string

var appendCmd = &cobra.Command{
	Use:   "append [ruta_remota] [contenido]",
	Short: "Agrega contenido al final de un archivo",
	Long: `Agrega texto al final de un archivo remoto existente.
Útil para agregar líneas a configuraciones o logs.

Ejemplos:
  sshcli append /etc/hosts "192.168.1.100 miservidor"
  sshcli append --server prod /var/www/.env "DEBUG=false"
  sshcli append /home/user/.bashrc "export PATH=$PATH:/opt/bin"`,
	Args: cobra.ExactArgs(2),
	RunE: runAppend,
}

func init() {
	rootCmd.AddCommand(appendCmd)
	appendCmd.Flags().StringVarP(&appendServer, "server", "s", "", "Servidor específico a usar")
}

func runAppend(cmd *cobra.Command, args []string) error {
	remotePath := args[0]
	content := args[1]

	client, _, err := getClient(appendServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	appendCommand := fmt.Sprintf("echo '%s' >> %s", content, remotePath)

	if _, err := client.Run(appendCommand); err != nil {
		return fmt.Errorf("error al agregar contenido: %v", err)
	}

	fmt.Printf("Contenido agregado a: %s\n", remotePath)
	return nil
}
