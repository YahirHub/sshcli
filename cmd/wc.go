package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	wcLines  bool
	wcWords  bool
	wcBytes  bool
	wcServer string
)

var wcCmd = &cobra.Command{
	Use:   "wc [archivo]",
	Short: "Cuenta líneas, palabras y bytes de un archivo",
	Long: `Cuenta líneas, palabras y bytes de un archivo remoto.
Por defecto muestra las tres métricas.

Ejemplos:
  sshcli wc /app/main.py
  sshcli wc -l /app/main.py              # Solo líneas
  sshcli wc --server prod /var/log/app.log

Casos de uso para agentes:
  - Medir tamaño de archivos de código
  - Verificar si archivo creció/decreció
  - Estimar complejidad por líneas
  - Monitorear tamaño de logs`,
	Args: cobra.ExactArgs(1),
	RunE: runWc,
}

func init() {
	rootCmd.AddCommand(wcCmd)
	wcCmd.Flags().BoolVarP(&wcLines, "lines", "l", false, "Mostrar solo líneas")
	wcCmd.Flags().BoolVarP(&wcWords, "words", "w", false, "Mostrar solo palabras")
	wcCmd.Flags().BoolVarP(&wcBytes, "bytes", "c", false, "Mostrar solo bytes")
	wcCmd.Flags().StringVarP(&wcServer, "server", "s", "", "Servidor específico a usar")
}

func runWc(cmd *cobra.Command, args []string) error {
	remotePath := args[0]

	client, _, err := getClient(wcServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	wcCommand := "wc"
	if wcLines {
		wcCommand += " -l"
	} else if wcWords {
		wcCommand += " -w"
	} else if wcBytes {
		wcCommand += " -c"
	}
	wcCommand += " " + remotePath

	output, err := client.Run(wcCommand)
	if err != nil {
		return fmt.Errorf("error al contar: %v", err)
	}

	fmt.Print(output)
	return nil
}
