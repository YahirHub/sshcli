package cmd

import (
	"fmt"
	"sshcli/internal/paths"

	"github.com/spf13/cobra"
)

var existsServer string

var existsCmd = &cobra.Command{
	Use:   "exists [ruta_remota]",
	Short: "Verifica si un archivo existe",
	Args:  cobra.ExactArgs(1),
	RunE:  runExists,
}

func init() {
	rootCmd.AddCommand(existsCmd)
	existsCmd.Flags().StringVarP(&existsServer, "server", "s", "", "Servidor específico")
}

func runExists(cmd *cobra.Command, args []string) error {
	remotePath := paths.ToRemote(args[0])

	client, _, err := getClient(existsServer)
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.Run(fmt.Sprintf("test -e '%s'", remotePath))
	if err != nil {
		fmt.Println("NO")
		return fmt.Errorf("no existe: %s", remotePath)
	}

	fmt.Println("SI")
	return nil
}