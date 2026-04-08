package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	chmodRecursive bool
	chmodServer    string
)

var chmodCmd = &cobra.Command{
	Use:   "chmod [permisos] [ruta_remota]",
	Short: "Cambia permisos de archivos o directorios",
	Long: `Cambia los permisos de un archivo o directorio remoto.
Usa notación octal (755, 644) o simbólica (+x, u+w).

Ejemplos:
  sshcli chmod 755 /home/user/script.sh
  sshcli chmod -r 644 /var/www/app
  sshcli chmod +x /home/user/run.sh`,
	Args: cobra.ExactArgs(2),
	RunE: runChmod,
}

func init() {
	rootCmd.AddCommand(chmodCmd)
	chmodCmd.Flags().BoolVarP(&chmodRecursive, "recursive", "r", false, "Aplicar recursivamente")
	chmodCmd.Flags().StringVarP(&chmodServer, "server", "s", "", "Servidor específico a usar")
}

func runChmod(cmd *cobra.Command, args []string) error {
	permissions := args[0]
	remotePath := args[1]

	client, _, err := getClient(chmodServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	chmodCommand := "chmod"
	if chmodRecursive {
		chmodCommand += " -R"
	}
	chmodCommand += fmt.Sprintf(" %s %s", permissions, remotePath)

	if _, err := client.Run(chmodCommand); err != nil {
		return fmt.Errorf("error al cambiar permisos: %v", err)
	}

	fmt.Printf("Permisos actualizados: %s -> %s\n", remotePath, permissions)
	return nil
}
