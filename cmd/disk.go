package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var diskServer string

var diskCmd = &cobra.Command{
	Use:   "disk [ruta]",
	Short: "Muestra uso de disco",
	Long: `Muestra el uso de disco del sistema o de una ruta específica.

Ejemplos:
  sshcli disk                         # Uso general del sistema
  sshcli disk /var/www                # Uso de directorio específico
  sshcli disk --server prod /app

Casos de uso para agentes:
  - Verificar espacio disponible
  - Identificar directorios grandes
  - Prevenir problemas de espacio`,
	Args: cobra.MaximumNArgs(1),
	RunE: runDisk,
}

func init() {
	rootCmd.AddCommand(diskCmd)
	diskCmd.Flags().StringVarP(&diskServer, "server", "s", "", "Servidor específico a usar")
}

func runDisk(cmd *cobra.Command, args []string) error {
	client, _, err := getClient(diskServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	var diskCommand string
	if len(args) > 0 {
		diskCommand = fmt.Sprintf("du -sh %s/* 2>/dev/null | sort -rh | head -20", args[0])
	} else {
		diskCommand = "df -h"
	}

	output, err := client.Run(diskCommand)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	fmt.Print(output)
	return nil
}
