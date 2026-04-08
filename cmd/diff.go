package cmd

import (
	"fmt"
	"os"
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
	Long: `Compara un archivo local con su versión en el servidor remoto.
Muestra las diferencias en formato unificado.

Ejemplos:
  sshcli diff ./main.py /app/main.py
  sshcli diff ./config.json /etc/app/config.json
  sshcli diff --server prod ./deploy.sh /opt/scripts/deploy.sh
  sshcli diff -c 5 ./app.py /app/app.py    # 5 líneas de contexto

Casos de uso para agentes:
  - Verificar cambios antes de subir
  - Comparar versiones local/remota
  - Detectar modificaciones no esperadas
  - Validar sincronización`,
	Args: cobra.ExactArgs(2),
	RunE: runDiff,
}

func init() {
	rootCmd.AddCommand(diffCmd)
	diffCmd.Flags().IntVarP(&diffContext, "context", "c", 3, "Líneas de contexto")
	diffCmd.Flags().StringVarP(&diffServer, "server", "s", "", "Servidor específico a usar")
}

func runDiff(cmd *cobra.Command, args []string) error {
	localPath := args[0]
	remotePath := args[1]

	// Leer archivo local
	localData, err := os.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("error al leer archivo local: %v", err)
	}

	client, _, err := getClient(diffServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	// Leer archivo remoto
	remoteData, err := client.ReadFile(remotePath)
	if err != nil {
		return fmt.Errorf("error al leer archivo remoto: %v", err)
	}

	localLines := strings.Split(string(localData), "\n")
	remoteLines := strings.Split(string(remoteData), "\n")

	// Comparación simple línea por línea
	if string(localData) == string(remoteData) {
		fmt.Println("Los archivos son idénticos")
		return nil
	}

	fmt.Printf("--- %s (local)\n", localPath)
	fmt.Printf("+++ %s (remoto)\n", remotePath)
	fmt.Println()

	maxLines := len(localLines)
	if len(remoteLines) > maxLines {
		maxLines = len(remoteLines)
	}

	inDiff := false
	diffStart := 0

	for i := 0; i < maxLines; i++ {
		var localLine, remoteLine string
		hasLocal := i < len(localLines)
		hasRemote := i < len(remoteLines)

		if hasLocal {
			localLine = localLines[i]
		}
		if hasRemote {
			remoteLine = remoteLines[i]
		}

		if localLine != remoteLine {
			if !inDiff {
				inDiff = true
				diffStart = i + 1
				fmt.Printf("@@ línea %d @@\n", diffStart)
			}
			if hasLocal {
				fmt.Printf("- %s\n", localLine)
			}
			if hasRemote {
				fmt.Printf("+ %s\n", remoteLine)
			}
		} else {
			if inDiff {
				inDiff = false
				fmt.Println()
			}
		}
	}

	return nil
}
