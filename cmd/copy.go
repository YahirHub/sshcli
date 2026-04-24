package cmd

import (
	"fmt"
	"sshcli/internal/paths"

	"github.com/spf13/cobra"
)

var copyServer string

var copyCmd = &cobra.Command{
	Use:   "copy [origen_remoto] [destino_remoto]",
	Short: "Copia archivos o directorios dentro del servidor",
	Long: `Copia archivos o directorios en el servidor remoto.
Ambas rutas son tratadas como rutas remotas de Linux y normalizadas automáticamente.

Ejemplos:
  sshcli copy /etc/nginx/nginx.conf /etc/nginx/nginx.conf.bak
  sshcli copy --server prod /var/www/app /var/www/app_backup`,
	Args: cobra.ExactArgs(2),
	RunE: runCopy,
}

func init() {
	rootCmd.AddCommand(copyCmd)
	copyCmd.Flags().StringVarP(&copyServer, "server", "s", "", "Servidor específico a usar")
}

func runCopy(cmd *cobra.Command, args []string) error {
	// 1. Normalizar rutas usando el motor global
	source := paths.ToRemote(args[0])
	dest := paths.ToRemote(args[1])

	// 2. Obtener cliente
	client, _, err := getClient(copyServer)
	if err != nil {
		return fmt.Errorf("error de conexión: %v", err)
	}
	defer client.Close()

	// 3. Ejecutar comando cp con rutas protegidas por comillas simples
	// -r permite copiar directorios recursivamente
	copyCommand := fmt.Sprintf("cp -r '%s' '%s'", source, dest)
	
	if _, err := client.Run(copyCommand); err != nil {
		return fmt.Errorf("error al copiar en el servidor: %v", err)
	}

	fmt.Printf("✓ Copiado exitosamente: %s -> %s\n", source, dest)
	return nil
}