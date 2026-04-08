package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	daemonName   string
	daemonLog    string
	daemonServer string
)

var daemonCmd = &cobra.Command{
	Use:   "daemon [comando]",
	Short: "Ejecuta un comando en segundo plano",
	Long: `Ejecuta un comando como proceso en segundo plano (daemon).
Opcionalmente guarda los logs en un archivo.

Ejemplos:
  sshcli daemon "python /app/server.py"
  sshcli daemon "node /app/index.js" --name myapp
  sshcli daemon "gunicorn app:app" --log /var/log/app.log
  sshcli daemon --server prod "python manage.py runserver"

Casos de uso para agentes:
  - Iniciar servidor de aplicación
  - Ejecutar trabajos de larga duración
  - Correr scripts en background`,
	Args: cobra.ExactArgs(1),
	RunE: runDaemon,
}

func init() {
	rootCmd.AddCommand(daemonCmd)
	daemonCmd.Flags().StringVarP(&daemonName, "name", "n", "", "Nombre identificador del daemon")
	daemonCmd.Flags().StringVarP(&daemonLog, "log", "l", "", "Archivo para guardar logs")
	daemonCmd.Flags().StringVarP(&daemonServer, "server", "s", "", "Servidor específico a usar")
}

func runDaemon(cmd *cobra.Command, args []string) error {
	command := args[0]

	client, _, err := getClient(daemonServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	var daemonCommand string
	if daemonLog != "" {
		daemonCommand = fmt.Sprintf("nohup %s > %s 2>&1 &", command, daemonLog)
	} else {
		daemonCommand = fmt.Sprintf("nohup %s > /dev/null 2>&1 &", command)
	}

	// Ejecutar y obtener PID
	fullCmd := fmt.Sprintf("%s echo $!", daemonCommand)
	output, err := client.Run(fullCmd)
	if err != nil {
		return fmt.Errorf("error al iniciar daemon: %v", err)
	}

	fmt.Printf("Daemon iniciado: %s\n", command)
	if daemonName != "" {
		fmt.Printf("Nombre: %s\n", daemonName)
	}
	if daemonLog != "" {
		fmt.Printf("Logs en: %s\n", daemonLog)
	}
	if output != "" {
		fmt.Printf("PID: %s", output)
	}

	return nil
}
