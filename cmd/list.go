package cmd

import (
	"fmt"
	"sshcli/internal/paths"

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
	Args:  cobra.MaximumNArgs(1),
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&listLong, "long", "l", false, "Formato largo")
	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "Incluir ocultos")
	listCmd.Flags().StringVarP(&listServer, "server", "s", "", "Servidor específico")
}

func runList(cmd *cobra.Command, args []string) error {
	remotePath := "/"
	if len(args) > 0 {
		remotePath = paths.ToRemote(args[0])
	}

	client, _, err := getClient(listServer)
	if err != nil {
		return err
	}
	defer client.Close()

	lsCmd := "ls"
	if listLong { lsCmd += " -l" }
	if listAll { lsCmd += " -a" }
	
	// Envolvemos en comillas para proteger contra espacios
	output, err := client.Run(fmt.Sprintf("%s '%s'", lsCmd, remotePath))
	if err != nil {
		return fmt.Errorf("error al listar: %v", err)
	}

	fmt.Print(output)
	return nil
}