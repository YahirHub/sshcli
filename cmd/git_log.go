package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	gitLogCount  int
	gitLogOneline bool
	gitLogServer string
)

var gitLogCmd = &cobra.Command{
	Use:   "git-log [directorio]",
	Short: "Muestra el historial de commits",
	Long: `Muestra el historial de commits del repositorio Git remoto.

Ejemplos:
  sshcli git-log /app/proyecto
  sshcli git-log -n 5 /app/proyecto          # Últimos 5 commits
  sshcli git-log --oneline /app/proyecto     # Formato compacto
  sshcli git-log --server prod /var/www/app

Casos de uso para agentes:
  - Ver historial de cambios
  - Identificar commits recientes
  - Obtener hash de commits para rollback`,
	Args: cobra.MaximumNArgs(1),
	RunE: runGitLog,
}

func init() {
	rootCmd.AddCommand(gitLogCmd)
	gitLogCmd.Flags().IntVarP(&gitLogCount, "number", "n", 10, "Número de commits a mostrar")
	gitLogCmd.Flags().BoolVar(&gitLogOneline, "oneline", false, "Formato de una línea por commit")
	gitLogCmd.Flags().StringVarP(&gitLogServer, "server", "s", "", "Servidor específico a usar")
}

func runGitLog(cmd *cobra.Command, args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	client, _, err := getClient(gitLogServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	gitCmd := fmt.Sprintf("cd %s && git log -n %d", dir, gitLogCount)
	if gitLogOneline {
		gitCmd += " --oneline"
	}

	output, err := client.Run(gitCmd)
	if err != nil {
		return fmt.Errorf("error al obtener log: %v", err)
	}

	fmt.Print(output)
	return nil
}
