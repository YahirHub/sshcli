package cmd

import (
	"fmt"
	"sshcli/internal/paths"

	"github.com/spf13/cobra"
)

var infoServer string

var infoCmd = &cobra.Command{
	Use:   "info [ruta_remota]",
	Short: "Muestra información detallada de un archivo o directorio",
	Args:  cobra.ExactArgs(1),
	RunE:  runInfo,
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringVarP(&infoServer, "server", "s", "", "Servidor específico a usar")
}

func runInfo(cmd *cobra.Command, args []string) error {
	remotePath := paths.ToRemote(args[0])

	client, _, err := getClient(infoServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	// Envolvemos en comillas simples para proteger contra espacios de Windows
	output, err := client.Run(fmt.Sprintf("stat '%s' 2>/dev/null || ls -la '%s'", remotePath, remotePath))
	if err != nil {
		return fmt.Errorf("error al obtener información: %v", err)
	}

	fmt.Print(output)
	return nil
}