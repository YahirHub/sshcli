package cmd

import (
	"fmt"
	"os"
	"strings"

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
	envFile := cleanRemotePath(args[0])
	keyValue := decodeEscapes(args[1])

	key, _, ok := strings.Cut(keyValue, "=")
	if !ok || strings.TrimSpace(key) == "" {
		return fmt.Errorf("formato inválido, usar: KEY=value")
	}
	key = strings.TrimSpace(key)

	client, _, err := getClient(envSetServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	var content string
	if client.FileExists(envFile) {
		data, err := client.ReadFile(envFile)
		if err != nil {
			return fmt.Errorf("error al leer .env: %v", err)
		}
		content = strings.ReplaceAll(string(data), "\r\n", "\n")
	}

	lines := []string{}
	if content != "" {
		lines = strings.Split(content, "\n")
	}

	updated := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, key+"=") {
			lines[i] = keyValue
			updated = true
		}
	}

	if !updated {
		for len(lines) > 0 && lines[len(lines)-1] == "" {
			lines = lines[:len(lines)-1]
		}
		lines = append(lines, keyValue)
	}

	newContent := strings.Join(lines, "\n")
	if newContent != "" && !strings.HasSuffix(newContent, "\n") {
		newContent += "\n"
	}

	if err := client.WriteFile(envFile, []byte(newContent), os.FileMode(0644)); err != nil {
		return fmt.Errorf("error al escribir .env: %v", err)
	}

	if updated {
		fmt.Printf("Actualizado: %s\n", key)
	} else {
		fmt.Printf("Agregado: %s\n", key)
	}
	return nil
}
