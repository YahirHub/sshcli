package cmd

import (
	"fmt"
	"sshcli/internal/paths"

	"github.com/spf13/cobra"
)

var (
	findName   string
	findType   string
	findServer string
)

var findCmd = &cobra.Command{
	Use:   "find [ruta_remota]",
	Short: "Busca archivos y directorios en el servidor remoto",
	Long: `Busca archivos y directorios usando patrones.
Limpia automáticamente la ruta de búsqueda para evitar conflictos con Git Bash.

Ejemplos:
  sshcli find /home/user -n "*.py"
  sshcli find /etc -t f -n "*.conf"`,
	Args: cobra.ExactArgs(1),
	RunE: runFind,
}

func init() {
	rootCmd.AddCommand(findCmd)
	findCmd.Flags().StringVarP(&findName, "name", "n", "", "Patrón de nombre a buscar")
	findCmd.Flags().StringVarP(&findType, "type", "t", "", "Tipo: f=archivo, d=directorio")
	findCmd.Flags().StringVarP(&findServer, "server", "s", "", "Servidor específico a usar")
}

func runFind(cmd *cobra.Command, args []string) error {
	// 1. Normalizar ruta
	remotePath := paths.ToRemote(args[0])

	// 2. Conectar
	client, _, err := getClient(findServer)
	if err != nil {
		return fmt.Errorf("error de conexión: %v", err)
	}
	defer client.Close()

	// 3. Construir comando find con comillas simples
	findCommand := fmt.Sprintf("find '%s'", remotePath)
	if findType != "" {
		findCommand += fmt.Sprintf(" -type %s", findType)
	}
	if findName != "" {
		findCommand += fmt.Sprintf(" -name '%s'", findName)
	}
	findCommand += " 2>/dev/null"

	output, err := client.Run(findCommand)
	if err != nil {
		return fmt.Errorf("error en búsqueda remota: %v", err)
	}

	if output == "" {
		fmt.Println("Sin resultados")
		return nil
	}

	fmt.Print(output)
	return nil
}