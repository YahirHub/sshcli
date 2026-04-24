package cmd

import (
	"fmt"
	"sshcli/internal/paths"

	"github.com/spf13/cobra"
)

var diskServer string

var diskCmd = &cobra.Command{
	Use:   "disk [ruta]",
	Short: "Muestra uso de disco",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runDisk,
}

func init() {
	rootCmd.AddCommand(diskCmd)
	diskCmd.Flags().StringVarP(&diskServer, "server", "s", "", "Servidor específico a usar")
}

func runDisk(cmd *cobra.Command, args []string) error {
	client, _, err := getClient(diskServer)
	if err != nil {
		return err
	}
	defer client.Close()

	var diskCommand string
	if len(args) > 0 {
		remotePath := paths.ToRemote(args[0])
		if remotePath == "/" {
			diskCommand = "df -h /"
		} else {
			diskCommand = fmt.Sprintf("du -sh %s/* 2>/dev/null | sort -rh | head -20", remotePath)
		}
	} else {
		diskCommand = "df -h"
	}

	output, err := client.Run(diskCommand)
	if err != nil {
		return fmt.Errorf("error al obtener disco: %v", err)
	}

	fmt.Print(output)
	return nil
}