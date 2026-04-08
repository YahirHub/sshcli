package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	tailLines  int
	tailFollow bool
	tailServer string
)

var tailCmd = &cobra.Command{
	Use:   "tail [archivo]",
	Short: "Muestra las últimas líneas de un archivo",
	Long: `Muestra las últimas líneas de un archivo remoto.
Ideal para ver logs y archivos que crecen continuamente.

Ejemplos:
  sshcli tail /var/log/app.log
  sshcli tail -n 100 /var/log/nginx/error.log
  sshcli tail --server prod /var/log/app.log
  sshcli tail -f /var/log/app.log    # Seguir en tiempo real (3 segundos)

Casos de uso para agentes:
  - Monitorear logs de aplicación
  - Ver errores recientes
  - Verificar output de procesos
  - Debug en tiempo real`,
	Args: cobra.ExactArgs(1),
	RunE: runTail,
}

func init() {
	rootCmd.AddCommand(tailCmd)
	tailCmd.Flags().IntVarP(&tailLines, "lines", "n", 20, "Número de líneas a mostrar")
	tailCmd.Flags().BoolVarP(&tailFollow, "follow", "f", false, "Seguir el archivo (muestra actualizaciones)")
	tailCmd.Flags().StringVarP(&tailServer, "server", "s", "", "Servidor específico a usar")
}

func runTail(cmd *cobra.Command, args []string) error {
	remotePath := args[0]

	client, _, err := getClient(tailServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	var tailCommand string
	if tailFollow {
		// Para follow, hacemos múltiples lecturas
		tailCommand = fmt.Sprintf("tail -n %d %s && sleep 2 && tail -n 5 %s", tailLines, remotePath, remotePath)
	} else {
		tailCommand = fmt.Sprintf("tail -n %s %s", strconv.Itoa(tailLines), remotePath)
	}

	output, err := client.Run(tailCommand)
	if err != nil {
		return fmt.Errorf("error al leer archivo: %v", err)
	}

	fmt.Print(output)
	return nil
}
