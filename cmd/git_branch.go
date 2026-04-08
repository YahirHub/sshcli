package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	gitBranchCreate string
	gitBranchDelete string
	gitBranchServer string
)

var gitBranchCmd = &cobra.Command{
	Use:   "git-branch [directorio]",
	Short: "Gestiona ramas del repositorio Git",
	Long: `Lista, crea o elimina ramas en el repositorio Git remoto.

Ejemplos:
  sshcli git-branch /app/proyecto                    # Listar ramas
  sshcli git-branch /app/proyecto --create feature   # Crear rama
  sshcli git-branch /app/proyecto --delete old-feat  # Eliminar rama
  sshcli git-branch --server prod /var/www/app

Casos de uso para agentes:
  - Ver ramas disponibles
  - Crear rama para nueva feature
  - Limpiar ramas obsoletas`,
	Args: cobra.MaximumNArgs(1),
	RunE: runGitBranch,
}

func init() {
	rootCmd.AddCommand(gitBranchCmd)
	gitBranchCmd.Flags().StringVarP(&gitBranchCreate, "create", "c", "", "Crear nueva rama")
	gitBranchCmd.Flags().StringVarP(&gitBranchDelete, "delete", "d", "", "Eliminar rama")
	gitBranchCmd.Flags().StringVarP(&gitBranchServer, "server", "s", "", "Servidor específico a usar")
}

func runGitBranch(cmd *cobra.Command, args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	client, _, err := getClient(gitBranchServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	var gitCmd string
	if gitBranchCreate != "" {
		gitCmd = fmt.Sprintf("cd %s && git checkout -b %s", dir, gitBranchCreate)
	} else if gitBranchDelete != "" {
		gitCmd = fmt.Sprintf("cd %s && git branch -d %s", dir, gitBranchDelete)
	} else {
		gitCmd = fmt.Sprintf("cd %s && git branch -a", dir)
	}

	output, err := client.Run(gitCmd)
	if err != nil {
		return fmt.Errorf("error: %v\n%s", err, output)
	}

	fmt.Print(output)
	return nil
}
