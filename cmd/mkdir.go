package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	mkdirParents bool
	mkdirServer  string
)

var mkdirCmd = &cobra.Command{
	Use:   "mkdir [ruta_remota]",
	Short: "Crea un directorio en el servidor remoto",
	Long: `Crea un directorio en el servidor remoto.
Usa -p para crear directorios padres si no existen.

Ejemplos:
  sshcli mkdir /home/user/nuevo_directorio
  sshcli mkdir -p --server dev /home/user/proyecto/src/components`,
	Args: cobra.ExactArgs(1),
	RunE: runMkdir,
}

func init() {
	rootCmd.AddCommand(mkdirCmd)
	mkdirCmd.Flags().BoolVarP(&mkdirParents, "parents", "p", false, "Crear directorios padres si no existen")
	mkdirCmd.Flags().StringVarP(&mkdirServer, "server", "s", "", "Servidor específico a usar")
}

func runMkdir(cmd *cobra.Command, args []string) error {
	remotePath := args[0]

	client, _, err := getClient(mkdirServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	mkdirCommand := "mkdir"
	if mkdirParents {
		mkdirCommand += " -p"
	}
	mkdirCommand += " " + remotePath

	if _, err := client.Run(mkdirCommand); err != nil {
		return fmt.Errorf("error al crear directorio: %v", err)
	}

	fmt.Printf("Directorio creado: %s\n", remotePath)
	return nil
}
