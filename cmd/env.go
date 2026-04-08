package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var envServer string

var envCmd = &cobra.Command{
	Use:   "env [archivo_env]",
	Short: "Muestra o gestiona variables de entorno",
	Long: `Muestra variables de entorno del sistema o de un archivo .env.
Sin argumentos muestra las variables del sistema.

Ejemplos:
  sshcli env                          # Variables del sistema
  sshcli env /app/.env                # Contenido de archivo .env
  sshcli env --server prod /var/www/.env

Casos de uso para agentes:
  - Verificar configuración de entorno
  - Leer valores de .env
  - Diagnosticar problemas de configuración`,
	Args: cobra.MaximumNArgs(1),
	RunE: runEnv,
}

func init() {
	rootCmd.AddCommand(envCmd)
	envCmd.Flags().StringVarP(&envServer, "server", "s", "", "Servidor específico a usar")
}

func runEnv(cmd *cobra.Command, args []string) error {
	client, _, err := getClient(envServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	var envCommand string
	if len(args) > 0 {
		envCommand = fmt.Sprintf("cat %s 2>/dev/null || echo 'Archivo no encontrado'", args[0])
	} else {
		envCommand = "env | sort"
	}

	output, err := client.Run(envCommand)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	fmt.Print(output)
	return nil
}
