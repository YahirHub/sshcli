package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	serviceAction string
	serviceServer string
)

var serviceCmd = &cobra.Command{
	Use:   "service [nombre] [accion]",
	Short: "Gestiona servicios del sistema (systemctl)",
	Long: `Gestiona servicios del sistema usando systemctl.
Acciones: start, stop, restart, status, enable, disable.

Ejemplos:
  sshcli service nginx status
  sshcli service nginx restart
  sshcli service postgresql start
  sshcli service --server prod nginx reload
  sshcli service myapp enable

Casos de uso para agentes:
  - Reiniciar servicios después de cambios
  - Verificar estado de servicios
  - Habilitar servicios al inicio`,
	Args: cobra.ExactArgs(2),
	RunE: runService,
}

func init() {
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.Flags().StringVarP(&serviceServer, "server", "s", "", "Servidor específico a usar")
}

func runService(cmd *cobra.Command, args []string) error {
	serviceName := args[0]
	action := args[1]

	validActions := map[string]bool{
		"start": true, "stop": true, "restart": true,
		"status": true, "enable": true, "disable": true,
		"reload": true,
	}

	if !validActions[action] {
		return fmt.Errorf("acción inválida: %s (usar: start, stop, restart, status, enable, disable, reload)", action)
	}

	client, _, err := getClient(serviceServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	serviceCommand := fmt.Sprintf("sudo systemctl %s %s", action, serviceName)
	output, err := client.Run(serviceCommand)
	if err != nil && action != "status" {
		return fmt.Errorf("error al ejecutar %s en %s: %v\n%s", action, serviceName, err, output)
	}

	if output != "" {
		fmt.Print(output)
	} else {
		fmt.Printf("Servicio %s: %s ejecutado\n", serviceName, action)
	}
	return nil
}
