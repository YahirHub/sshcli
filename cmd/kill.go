package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	killSignal string
	killServer string
)

var killCmd = &cobra.Command{
	Use:   "kill [pid_o_nombre]",
	Short: "Termina un proceso por PID o nombre",
	Long: `Termina un proceso en el servidor remoto.
Puede especificar PID numérico o nombre del proceso.

Ejemplos:
  sshcli kill 12345                   # Por PID
  sshcli kill python                  # Por nombre (pkill)
  sshcli kill --signal 9 12345        # SIGKILL
  sshcli kill --signal HUP nginx      # Reload nginx
  sshcli kill --server prod gunicorn

Señales comunes:
  15 (TERM) - Terminación graceful (default)
  9  (KILL) - Terminación forzada
  1  (HUP)  - Reload configuración
  2  (INT)  - Interrupción (Ctrl+C)`,
	Args: cobra.ExactArgs(1),
	RunE: runKill,
}

func init() {
	rootCmd.AddCommand(killCmd)
	killCmd.Flags().StringVar(&killSignal, "signal", "15", "Señal a enviar (default: 15/TERM)")
	killCmd.Flags().StringVarP(&killServer, "server", "s", "", "Servidor específico a usar")
}

func runKill(cmd *cobra.Command, args []string) error {
	target := args[0]

	client, _, err := getClient(killServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	// Determinar si es PID o nombre
	var killCommand string
	if isNumeric(target) {
		killCommand = fmt.Sprintf("kill -%s %s", killSignal, target)
	} else {
		killCommand = fmt.Sprintf("pkill -%s %s", killSignal, target)
	}

	output, err := client.Run(killCommand)
	if err != nil {
		return fmt.Errorf("error al terminar proceso: %v", err)
	}

	if output != "" {
		fmt.Print(output)
	}
	fmt.Printf("Señal %s enviada a: %s\n", killSignal, target)
	return nil
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}
