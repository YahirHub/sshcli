package cmd

import (
	"fmt"
	"sshcli/internal/paths"

	"github.com/spf13/cobra"
)

var (
	tailLines  int
	tailFollow bool
	tailServer string
)

var tailCmd = &cobra.Command{
	Use:   "tail [archivo]",
	Short: "Muestra las últimas líneas de un archivo",
	Args:  cobra.ExactArgs(1),
	RunE:  runTail,
}

func init() {
	rootCmd.AddCommand(tailCmd)
	tailCmd.Flags().IntVarP(&tailLines, "lines", "n", 20, "Número de líneas a mostrar")
	tailCmd.Flags().BoolVarP(&tailFollow, "follow", "f", false, "Seguir el archivo")
	tailCmd.Flags().StringVarP(&tailServer, "server", "s", "", "Servidor específico a usar")
}

func runTail(cmd *cobra.Command, args []string) error {
	remotePath := paths.ToRemote(args[0])

	client, _, err := getClient(tailServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	tailCommand := "tail"
	if tailFollow {
		tailCommand += " -f"
	}
	tailCommand += fmt.Sprintf(" -n %d", tailLines)

	output, err := client.Run(fmt.Sprintf("%s '%s'", tailCommand, remotePath))
	if err != nil {
		return fmt.Errorf("error al leer archivo: %v", err)
	}

	fmt.Print(output)
	return nil
}