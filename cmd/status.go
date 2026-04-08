package cmd

import (
	"fmt"

	"sshcli/internal/config"

	"github.com/spf13/cobra"
)

var statusServer string

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Muestra el estado de la conexión y configuración",
	Long: `Muestra información sobre la configuración actual
y verifica el estado de la conexión al servidor remoto.

Ejemplos:
  sshcli status
  sshcli status --server produccion`,
	RunE: runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().StringVarP(&statusServer, "server", "s", "", "Servidor específico a verificar")
}

func runStatus(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Estado: NO CONFIGURADO")
		fmt.Println("Ejecuta 'sshcli server add' para agregar un servidor")
		return nil
	}

	if len(cfg.Servers) == 0 {
		fmt.Println("Estado: SIN SERVIDORES")
		fmt.Println("Ejecuta 'sshcli server add' para agregar uno")
		return nil
	}

	client, server, err := getClient(statusServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	fmt.Println("=== Configuración ===")
	fmt.Printf("Servidor: %s\n", server.Name)
	fmt.Printf("Host: %s\n", server.Host)
	fmt.Printf("Puerto: %d\n", server.Port)
	fmt.Printf("Usuario: %s\n", server.User)
	fmt.Println()

	fmt.Print("Verificando conexión... ")

	output, err := client.Run("hostname && uptime")
	if err != nil {
		fmt.Println("FALLO")
		return fmt.Errorf("error al verificar: %v", err)
	}

	fmt.Println("OK")
	fmt.Println()
	fmt.Println("=== Información del servidor ===")
	fmt.Print(output)

	return nil
}
