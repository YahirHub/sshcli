package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	portsListen bool
	portsServer string
)

var portsCmd = &cobra.Command{
	Use:   "ports",
	Short: "Lista puertos en uso o escuchando",
	Long: `Muestra los puertos de red en uso en el servidor remoto.
Por defecto muestra puertos escuchando (LISTEN).

Ejemplos:
  sshcli ports
  sshcli ports --listen              # Solo puertos escuchando
  sshcli ports --server prod

Casos de uso para agentes:
  - Verificar si puerto está disponible
  - Encontrar qué proceso usa un puerto
  - Diagnosticar conflictos de puertos`,
	RunE: runPorts,
}

func init() {
	rootCmd.AddCommand(portsCmd)
	portsCmd.Flags().BoolVarP(&portsListen, "listen", "l", true, "Mostrar solo puertos escuchando")
	portsCmd.Flags().StringVarP(&portsServer, "server", "s", "", "Servidor específico a usar")
}

func runPorts(cmd *cobra.Command, args []string) error {
	client, _, err := getClient(portsServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	var portsCommand string
	if portsListen {
		portsCommand = "ss -tlnp 2>/dev/null || netstat -tlnp 2>/dev/null"
	} else {
		portsCommand = "ss -tanp 2>/dev/null || netstat -tanp 2>/dev/null"
	}

	output, err := client.Run(portsCommand)
	if err != nil {
		return fmt.Errorf("error al listar puertos: %v", err)
	}

	fmt.Print(output)
	return nil
}
