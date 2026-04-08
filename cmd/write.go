package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var writeServer string

var writeCmd = &cobra.Command{
	Use:   "write [ruta_remota]",
	Short: "Escribe contenido desde stdin a un archivo remoto",
	Long: `Escribe contenido directamente a un archivo remoto.
Lee el contenido desde stdin, ideal para agentes de IA que necesitan
crear archivos con contenido específico.

Ejemplos:
  echo "contenido" | sshcli write /home/user/archivo.txt
  cat script.py | sshcli write --server dev /home/user/script.py
  sshcli write /home/user/config.json < config.json`,
	Args: cobra.ExactArgs(1),
	RunE: runWrite,
}

func init() {
	rootCmd.AddCommand(writeCmd)
	writeCmd.Flags().StringVarP(&writeServer, "server", "s", "", "Servidor específico a usar")
}

func runWrite(cmd *cobra.Command, args []string) error {
	remotePath := args[0]

	client, _, err := getClient(writeServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	content, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("error al leer stdin: %v", err)
	}

	dir := filepath.Dir(remotePath)
	if _, err := client.Run(fmt.Sprintf("mkdir -p %s", dir)); err != nil {
		return fmt.Errorf("error al crear directorio: %v", err)
	}

	if err := client.WriteFile(remotePath, content, 0644); err != nil {
		return fmt.Errorf("error al escribir archivo: %v", err)
	}

	fmt.Printf("Archivo creado: %s (%d bytes)\n", remotePath, len(content))
	return nil
}
