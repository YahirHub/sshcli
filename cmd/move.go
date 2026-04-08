package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var moveServer string

var moveCmd = &cobra.Command{
	Use:   "move [origen_remoto] [destino_remoto]",
	Short: "Mueve o renombra archivos o directorios",
	Long: `Mueve o renombra archivos y directorios en el servidor remoto.

Ejemplos:
  sshcli move /home/user/viejo.txt /home/user/nuevo.txt
  sshcli move --server prod /var/www/old_app /var/www/app`,
	Args: cobra.ExactArgs(2),
	RunE: runMove,
}

func init() {
	rootCmd.AddCommand(moveCmd)
	moveCmd.Flags().StringVarP(&moveServer, "server", "s", "", "Servidor específico a usar")
}

func runMove(cmd *cobra.Command, args []string) error {
	source := args[0]
	dest := args[1]

	client, _, err := getClient(moveServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	if _, err := client.Run(fmt.Sprintf("mv %s %s", source, dest)); err != nil {
		return fmt.Errorf("error al mover: %v", err)
	}

	fmt.Printf("Movido: %s -> %s\n", source, dest)
	return nil
}
