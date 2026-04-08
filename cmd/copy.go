package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var copyServer string

var copyCmd = &cobra.Command{
	Use:   "copy [origen_remoto] [destino_remoto]",
	Short: "Copia archivos o directorios dentro del servidor",
	Long: `Copia archivos o directorios en el servidor remoto.
El origen y destino son rutas remotas.

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
	source := args[0]
	dest := args[1]

	client, _, err := getClient(copyServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	if _, err := client.Run(fmt.Sprintf("cp -r %s %s", source, dest)); err != nil {
		return fmt.Errorf("error al copiar: %v", err)
	}

	fmt.Printf("Copiado: %s -> %s\n", source, dest)
	return nil
}
