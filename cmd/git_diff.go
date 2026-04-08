package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	gitDiffStaged bool
	gitDiffServer string
)

var gitDiffCmd = &cobra.Command{
	Use:   "git-diff [directorio]",
	Short: "Muestra los cambios en el repositorio Git",
	Long: `Muestra las diferencias en el repositorio Git remoto.
Por defecto muestra cambios no staged.

Ejemplos:
  sshcli git-diff /app/proyecto
  sshcli git-diff --staged /app/proyecto    # Cambios staged
  sshcli git-diff --server prod /var/www/app

Casos de uso para agentes:
  - Revisar cambios antes de commit
  - Verificar modificaciones realizadas
  - Inspeccionar código cambiado`,
	Args: cobra.MaximumNArgs(1),
	RunE: runGitDiff,
}

func init() {
	rootCmd.AddCommand(gitDiffCmd)
	gitDiffCmd.Flags().BoolVar(&gitDiffStaged, "staged", false, "Mostrar cambios staged")
	gitDiffCmd.Flags().StringVarP(&gitDiffServer, "server", "s", "", "Servidor específico a usar")
}

func runGitDiff(cmd *cobra.Command, args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	client, _, err := getClient(gitDiffServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	gitCmd := fmt.Sprintf("cd %s && git diff", dir)
	if gitDiffStaged {
		gitCmd += " --staged"
	}

	output, err := client.Run(gitCmd)
	if err != nil {
		return fmt.Errorf("error al obtener diff: %v", err)
	}

	if output == "" {
		fmt.Println("Sin cambios")
		return nil
	}

	fmt.Print(output)
	return nil
}
