package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	removeRecursive bool
	removeForce     bool
	removeServer    string
)

var removeCmd = &cobra.Command{
	Use:   "remove [ruta_remota]",
	Short: "Elimina un archivo o directorio del servidor remoto",
	Long: `Elimina un archivo o directorio del servidor remoto.
Usa -r para eliminar directorios recursivamente.
Usa -f para forzar la eliminación sin confirmación.

Ejemplos:
  sshcli remove /home/user/archivo.txt
  sshcli remove -r --server prod /home/user/carpeta
  sshcli remove -rf /tmp/temporal`,
	Args: cobra.ExactArgs(1),
	RunE: runRemove,
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolVarP(&removeRecursive, "recursive", "r", false, "Eliminar directorios recursivamente")
	removeCmd.Flags().BoolVarP(&removeForce, "force", "f", false, "Forzar eliminación sin confirmación")
	removeCmd.Flags().StringVarP(&removeServer, "server", "s", "", "Servidor específico a usar")
}

func runRemove(cmd *cobra.Command, args []string) error {
	remotePath := args[0]

	client, _, err := getClient(removeServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	rmCommand := "rm"
	if removeRecursive {
		rmCommand += " -r"
	}
	if removeForce {
		rmCommand += " -f"
	}
	rmCommand += " " + remotePath

	if _, err := client.Run(rmCommand); err != nil {
		return fmt.Errorf("error al eliminar: %v", err)
	}

	fmt.Printf("Eliminado: %s\n", remotePath)
	return nil
}
