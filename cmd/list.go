package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	listLong   bool
	listAll    bool
	listServer string
)

var listCmd = &cobra.Command{
	Use:   "list [ruta_remota]",
	Short: "Lista archivos y directorios en una ruta remota",
	Long: `Lista el contenido de un directorio remoto.
Por defecto lista el directorio home del usuario.

Ejemplos:
  sshcli list
  sshcli list /var/log
  sshcli list -l --server produccion /home/user
  sshcli list -a /etc`,
	Args: cobra.MaximumNArgs(1),
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&listLong, "long", "l", false, "Formato largo con detalles")
	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "Incluir archivos ocultos")
	listCmd.Flags().StringVarP(&listServer, "server", "s", "", "Servidor específico a usar")
}

func runList(cmd *cobra.Command, args []string) error {
	remotePath := "."
	if len(args) > 0 {
		remotePath = args[0]
	}

	client, _, err := getClient(listServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	lsCmd := "ls"
	if listLong {
		lsCmd += " -l"
	}
	if listAll {
		lsCmd += " -a"
	}
	lsCmd += " " + remotePath

	output, err := client.Run(lsCmd)
	if err != nil {
		return fmt.Errorf("error al listar directorio: %v", err)
	}

	fmt.Print(output)
	return nil
}
