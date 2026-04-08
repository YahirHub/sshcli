package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	envSetServer string
)

var envSetCmd = &cobra.Command{
	Use:   "env-set [archivo_env] [KEY=value]",
	Short: "Establece una variable en archivo .env",
	Long: `Agrega o actualiza una variable en un archivo .env.
Si la variable existe, la actualiza. Si no existe, la agrega.

Ejemplos:
  sshcli env-set /app/.env "DATABASE_URL=postgres://localhost/db"
  sshcli env-set /app/.env "DEBUG=false"
  sshcli env-set --server prod /var/www/.env "API_KEY=xxx123"

Casos de uso para agentes:
  - Configurar variables de entorno
  - Actualizar configuración de aplicación
  - Cambiar entre ambientes`,
	Args: cobra.ExactArgs(2),
	RunE: runEnvSet,
}

func init() {
	rootCmd.AddCommand(envSetCmd)
	envSetCmd.Flags().StringVarP(&envSetServer, "server", "s", "", "Servidor específico a usar")
}

func runEnvSet(cmd *cobra.Command, args []string) error {
	envFile := args[0]
	keyValue := args[1]

	client, _, err := getClient(envSetServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	// Extraer KEY del KEY=value
	var key string
	for i, c := range keyValue {
		if c == '=' {
			key = keyValue[:i]
			break
		}
	}

	if key == "" {
		return fmt.Errorf("formato inválido, usar: KEY=value")
	}

	// Comando para actualizar o agregar
	// Primero intenta actualizar, si no existe, agrega
	envCommand := fmt.Sprintf(`
		if grep -q "^%s=" %s 2>/dev/null; then
			sed -i "s|^%s=.*|%s|" %s
			echo "Actualizado: %s"
		else
			echo "%s" >> %s
			echo "Agregado: %s"
		fi
	`, key, envFile, key, keyValue, envFile, key, keyValue, envFile, key)

	output, err := client.Run(envCommand)
	if err != nil {
		return fmt.Errorf("error al establecer variable: %v", err)
	}

	fmt.Print(output)
	return nil
}
