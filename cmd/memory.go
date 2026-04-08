package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var memoryServer string

var memoryCmd = &cobra.Command{
	Use:   "memory",
	Short: "Muestra uso de memoria del sistema",
	Long: `Muestra el uso de memoria RAM del servidor remoto.

Ejemplos:
  sshcli memory
  sshcli memory --server prod

Casos de uso para agentes:
  - Monitorear recursos del servidor
  - Diagnosticar problemas de memoria
  - Verificar disponibilidad antes de deploy`,
	RunE: runMemory,
}

func init() {
	rootCmd.AddCommand(memoryCmd)
	memoryCmd.Flags().StringVarP(&memoryServer, "server", "s", "", "Servidor específico a usar")
}

func runMemory(cmd *cobra.Command, args []string) error {
	client, _, err := getClient(memoryServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	memCommand := "free -h"
	output, err := client.Run(memCommand)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	fmt.Print(output)
	return nil
}
