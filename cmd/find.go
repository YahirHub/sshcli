package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	findName   string
	findType   string
	findServer string
)

var findCmd = &cobra.Command{
	Use:   "find [ruta_remota]",
	Short: "Busca archivos y directorios en el servidor remoto",
	Long: `Busca archivos y directorios usando patrones.
Útil para agentes que necesitan localizar archivos específicos.

Ejemplos:
  sshcli find /home/user -name "*.py"
  sshcli find --server prod /var/www -name "config*"
  sshcli find /etc -type f -name "*.conf"`,
	Args: cobra.ExactArgs(1),
	RunE: runFind,
}

func init() {
	rootCmd.AddCommand(findCmd)
	findCmd.Flags().StringVarP(&findName, "name", "n", "", "Patrón de nombre a buscar")
	findCmd.Flags().StringVarP(&findType, "type", "t", "", "Tipo: f=archivo, d=directorio")
	findCmd.Flags().StringVarP(&findServer, "server", "s", "", "Servidor específico a usar")
}

func runFind(cmd *cobra.Command, args []string) error {
	remotePath := args[0]

	client, _, err := getClient(findServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	findCommand := fmt.Sprintf("find %s", remotePath)
	if findType != "" {
		findCommand += fmt.Sprintf(" -type %s", findType)
	}
	if findName != "" {
		findCommand += fmt.Sprintf(" -name '%s'", findName)
	}
	findCommand += " 2>/dev/null"

	output, err := client.Run(findCommand)
	if err != nil {
		return fmt.Errorf("error en búsqueda: %v", err)
	}

	if output == "" {
		fmt.Println("Sin resultados")
		return nil
	}

	fmt.Print(output)
	return nil
}
