package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	psFilter string
	psAll    bool
	psServer string
)

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "Lista procesos en ejecución",
	Long: `Lista los procesos en ejecución en el servidor remoto.
Puede filtrar por nombre de proceso.

Ejemplos:
  sshcli ps
  sshcli ps --filter python           # Solo procesos python
  sshcli ps --filter nginx
  sshcli ps --all                     # Todos los procesos
  sshcli ps --server prod

Casos de uso para agentes:
  - Verificar si aplicación está corriendo
  - Encontrar PID de procesos
  - Monitorear servicios`,
	RunE: runPs,
}

func init() {
	rootCmd.AddCommand(psCmd)
	psCmd.Flags().StringVarP(&psFilter, "filter", "f", "", "Filtrar por nombre de proceso")
	psCmd.Flags().BoolVarP(&psAll, "all", "a", false, "Mostrar todos los procesos")
	psCmd.Flags().StringVarP(&psServer, "server", "s", "", "Servidor específico a usar")
}

func runPs(cmd *cobra.Command, args []string) error {
	client, _, err := getClient(psServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	var psCommand string
	if psFilter != "" {
		psCommand = fmt.Sprintf("ps aux | grep -i %s | grep -v grep", psFilter)
	} else if psAll {
		psCommand = "ps aux"
	} else {
		psCommand = "ps aux | head -20"
	}

	output, err := client.Run(psCommand)
	if err != nil {
		if psFilter != "" {
			fmt.Printf("No se encontraron procesos con '%s'\n", psFilter)
			return nil
		}
		return fmt.Errorf("error: %v", err)
	}

	fmt.Print(output)
	return nil
}
