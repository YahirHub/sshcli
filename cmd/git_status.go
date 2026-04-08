package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var gitStatusServer string

var gitStatusCmd = &cobra.Command{
	Use:   "git-status [directorio]",
	Short: "Muestra el estado del repositorio Git",
	Long: `Muestra el estado actual del repositorio Git remoto.
Incluye archivos modificados, staged, untracked, etc.

Ejemplos:
  sshcli git-status /app/proyecto
  sshcli git-status --server prod /var/www/app
  sshcli git-status .

Casos de uso para agentes:
  - Verificar cambios pendientes antes de commit
  - Detectar archivos modificados
  - Confirmar estado limpio del repo`,
	Args: cobra.MaximumNArgs(1),
	RunE: runGitStatus,
}

func init() {
	rootCmd.AddCommand(gitStatusCmd)
	gitStatusCmd.Flags().StringVarP(&gitStatusServer, "server", "s", "", "Servidor específico a usar")
}

func runGitStatus(cmd *cobra.Command, args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	client, _, err := getClient(gitStatusServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	gitCmd := fmt.Sprintf("cd %s && git status", dir)
	output, err := client.Run(gitCmd)
	if err != nil {
		return fmt.Errorf("error al obtener estado: %v", err)
	}

	fmt.Print(output)
	return nil
}
