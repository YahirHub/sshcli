package cmd

import (
	"fmt"
	"sshcli/internal/paths"

	"github.com/spf13/cobra"
)

var readServer string

var readCmd = &cobra.Command{
	Use:   "read [ruta_remota]",
	Short: "Lee el contenido de un archivo remoto",
	Args:  cobra.ExactArgs(1),
	RunE:  runRead,
}

func init() {
	rootCmd.AddCommand(readCmd)
	readCmd.Flags().StringVarP(&readServer, "server", "s", "", "Servidor específico a usar")
}

func runRead(cmd *cobra.Command, args []string) error {
	remotePath := paths.ToRemote(args[0])

	client, _, err := getClient(readServer)
	if err != nil {
		return fmt.Errorf("error de conexión: %v", err)
	}
	defer client.Close()

	data, err := client.ReadFile(remotePath)
	if err != nil {
		return fmt.Errorf("error al leer archivo remoto %s: %v", remotePath, err)
	}

	fmt.Print(string(data))
	return nil
}