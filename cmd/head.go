package cmd

import (
	"fmt"
	"sshcli/internal/paths"
	"github.com/spf13/cobra"
)

var (
	headLines  int
	headServer string
)

var headCmd = &cobra.Command{
	Use:   "head [archivo]",
	Short: "Muestra las primeras líneas de un archivo",
	Args:  cobra.ExactArgs(1),
	RunE:  runHead,
}

func init() {
	rootCmd.AddCommand(headCmd)
	headCmd.Flags().IntVarP(&headLines, "lines", "n", 20, "Número de líneas a mostrar")
	headCmd.Flags().StringVarP(&headServer, "server", "s", "", "Servidor específico a usar")
}

func runHead(cmd *cobra.Command, args []string) error {
	remotePath := paths.ToRemote(args[0])

	client, _, err := getClient(headServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	output, err := client.Run(fmt.Sprintf("head -n %d '%s'", headLines, remotePath))
	if err != nil {
		return fmt.Errorf("error al leer archivo: %v", err)
	}

	fmt.Print(output)
	return nil
}