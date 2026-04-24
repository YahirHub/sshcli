package cmd

import (
	"fmt"
	"sshcli/internal/paths"

	"github.com/spf13/cobra"
)

var moveServer string

var moveCmd = &cobra.Command{
	Use:   "move [origen_remoto] [destino_remoto]",
	Short: "Mueve o renombra archivos o directorios",
	Long: `Mueve o renombra archivos y directorios en el servidor remoto.
Ambas rutas son normalizadas para evitar errores de shells de Windows.

Ejemplos:
  sshcli move /home/user/viejo.txt /home/user/nuevo.txt
  sshcli move --server prod /var/www/old_app /var/www/app`,
	Args: cobra.ExactArgs(2),
	RunE: runMove,
}

func init() {
	rootCmd.AddCommand(moveCmd)
	moveCmd.Flags().StringVarP(&moveServer, "server", "s", "", "Servidor específico a usar")
}

func runMove(cmd *cobra.Command, args []string) error {
	// 1. Normalizar rutas usando el motor global
	source := paths.ToRemote(args[0])
	dest := paths.ToRemote(args[1])

	// 2. Obtener cliente
	client, _, err := getClient(moveServer)
	if err != nil {
		return fmt.Errorf("error de conexión: %v", err)
	}
	defer client.Close()

	// 3. Ejecutar comando mv con rutas protegidas por comillas simples
	moveCommand := fmt.Sprintf("mv '%s' '%s'", source, dest)

	if _, err := client.Run(moveCommand); err != nil {
		return fmt.Errorf("error al mover en el servidor: %v", err)
	}

	fmt.Printf("✓ Movido exitosamente: %s -> %s\n", source, dest)
	return nil
}