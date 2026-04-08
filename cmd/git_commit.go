package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	gitCommitMsg    string
	gitCommitAll    bool
	gitCommitServer string
)

var gitCommitCmd = &cobra.Command{
	Use:   "git-commit [directorio]",
	Short: "Crea un commit en el repositorio Git",
	Long: `Crea un commit con los cambios staged en el repositorio Git remoto.
Usa -a para agregar y commitear todos los archivos modificados.

Ejemplos:
  sshcli git-commit /app/proyecto -m "fix: corregido bug en login"
  sshcli git-commit -a /app/proyecto -m "feat: nueva funcionalidad"
  sshcli git-commit --server prod /var/www/app -m "deploy: versión 2.0"

Casos de uso para agentes:
  - Guardar cambios realizados
  - Crear puntos de restauración
  - Documentar modificaciones`,
	Args: cobra.MaximumNArgs(1),
	RunE: runGitCommit,
}

func init() {
	rootCmd.AddCommand(gitCommitCmd)
	gitCommitCmd.Flags().StringVarP(&gitCommitMsg, "message", "m", "", "Mensaje del commit (requerido)")
	gitCommitCmd.Flags().BoolVarP(&gitCommitAll, "all", "a", false, "Agregar todos los archivos modificados")
	gitCommitCmd.Flags().StringVarP(&gitCommitServer, "server", "s", "", "Servidor específico a usar")
	gitCommitCmd.MarkFlagRequired("message")
}

func runGitCommit(cmd *cobra.Command, args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	client, _, err := getClient(gitCommitServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	var gitCmd string
	if gitCommitAll {
		gitCmd = fmt.Sprintf("cd %s && git add -A && git commit -m '%s'", dir, gitCommitMsg)
	} else {
		gitCmd = fmt.Sprintf("cd %s && git commit -m '%s'", dir, gitCommitMsg)
	}

	output, err := client.Run(gitCmd)
	if err != nil {
		return fmt.Errorf("error al hacer commit: %v\n%s", err, output)
	}

	fmt.Print(output)
	return nil
}
