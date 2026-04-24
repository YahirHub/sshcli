package cmd

import (
	"fmt"
	"os"
	"sshcli/internal/paths"
	"strings"

	"github.com/spf13/cobra"
)

var (
	diffContext int
	diffServer  string
)

var diffCmd = &cobra.Command{
	Use:   "diff [archivo_local] [archivo_remoto]",
	Short: "Compara un archivo local con uno remoto",
	Args:  cobra.ExactArgs(2),
	RunE:  runDiff,
}

func init() {
	rootCmd.AddCommand(diffCmd)
	diffCmd.Flags().IntVarP(&diffContext, "context", "c", 3, "Líneas de contexto")
	diffCmd.Flags().StringVarP(&diffServer, "server", "s", "", "Servidor específico a usar")
}

func runDiff(cmd *cobra.Command, args []string) error {
	localPath := paths.ToLocal(args[0])
	remotePath := paths.ToRemote(args[1])

	localData, err := os.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("error al leer archivo local: %v", err)
	}

	client, _, err := getClient(diffServer)
	if err != nil {
		return err
	}
	defer client.Close()

	remoteData, err := client.ReadFile(remotePath)
	if err != nil {
		return fmt.Errorf("error al leer archivo remoto: %v", err)
	}

	if string(localData) == string(remoteData) {
		fmt.Println("Los archivos son idénticos")
		return nil
	}

	// Lógica de diff simplificada para brevedad
	fmt.Printf("--- %s (local)\n+++ %s (remoto)\n", localPath, remotePath)
	localLines := strings.Split(string(localData), "\n")
	remoteLines := strings.Split(string(remoteData), "\n")
	
	for i := 0; i < len(localLines) || i < len(remoteLines); i++ {
		l := ""
		if i < len(localLines) { l = localLines[i] }
		r := ""
		if i < len(remoteLines) { r = remoteLines[i] }
		if l != r {
			fmt.Printf("- %s\n+ %s\n", l, r)
		}
	}

	return nil
}