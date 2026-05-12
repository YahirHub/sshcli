package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var appendServer string

var appendCmd = &cobra.Command{
	Use:   "append [ruta_remota] [contenido]",
	Short: "Agrega contenido al final de un archivo",
	Long:  `Agrega texto al final de un archivo remoto existente.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runAppend,
}

func init() {
	rootCmd.AddCommand(appendCmd)
	appendCmd.Flags().StringVarP(&appendServer, "server", "s", "", "Servidor específico a usar")
}

func runAppend(cmd *cobra.Command, args []string) error {
	remotePath := cleanRemotePath(args[0])
	content := []byte(decodeEscapes(args[1]))

	client, _, err := getClient(appendServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	existing, err := client.ReadFile(remotePath)
	if err != nil {
		return fmt.Errorf("error al leer archivo para append: %v", err)
	}

	newContent := append(existing, content...)
	if err := client.WriteFile(remotePath, newContent, os.FileMode(0644)); err != nil {
		return fmt.Errorf("error al agregar contenido: %v", err)
	}

	fmt.Printf("Contenido agregado a: %s\n", remotePath)
	return nil
}
