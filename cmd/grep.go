package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	grepRecursive bool
	grepIgnore    bool
	grepServer    string
)

var grepCmd = &cobra.Command{
	Use:   "grep [patron] [ruta_remota]",
	Short: "Busca texto dentro de archivos remotos",
	Long: `Busca un patrón de texto dentro de archivos remotos.
Útil para agentes que necesitan encontrar código o configuraciones.

Ejemplos:
  sshcli grep "import os" /home/user/script.py
  sshcli grep -r "DATABASE_URL" /var/www/app
  sshcli grep -i "error" /var/log/app.log`,
	Args: cobra.ExactArgs(2),
	RunE: runGrep,
}

func init() {
	rootCmd.AddCommand(grepCmd)
	grepCmd.Flags().BoolVarP(&grepRecursive, "recursive", "r", false, "Buscar recursivamente en directorios")
	grepCmd.Flags().BoolVarP(&grepIgnore, "ignore-case", "i", false, "Ignorar mayúsculas/minúsculas")
	grepCmd.Flags().StringVarP(&grepServer, "server", "s", "", "Servidor específico a usar")
}

func runGrep(cmd *cobra.Command, args[]string) error {
	pattern := args[0]
	remotePath := cleanRemotePath(args[1])

	client, _, err := getClient(grepServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	grepCommand := "grep"
	if grepRecursive {
		grepCommand += " -r"
	}
	if grepIgnore {
		grepCommand += " -i"
	}
	grepCommand += fmt.Sprintf(" -n '%s' '%s' 2>/dev/null", pattern, remotePath)

	output, err := client.Run(grepCommand)
	if err != nil {
		fmt.Println("Sin coincidencias")
		return nil
	}

	if output == "" {
		fmt.Println("Sin coincidencias")
		return nil
	}

	fmt.Print(output)
	return nil
}