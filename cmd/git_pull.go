package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var gitPullServer string

var gitPullCmd = &cobra.Command{
	Use:   "git-pull [directorio]",
	Short: "Actualiza el repositorio desde el remoto",
	Long: `Ejecuta git pull para actualizar el repositorio con cambios remotos.

Ejemplos:
  sshcli git-pull /app/proyecto
  sshcli git-pull --server prod /var/www/app

Casos de uso para agentes:
  - Sincronizar código con origen
  - Obtener últimos cambios del equipo
  - Actualizar antes de hacer deploy`,
	Args: cobra.MaximumNArgs(1),
	RunE: runGitPull,
}

func init() {
	rootCmd.AddCommand(gitPullCmd)
	gitPullCmd.Flags().StringVarP(&gitPullServer, "server", "s", "", "Servidor específico a usar")
}

func runGitPull(cmd *cobra.Command, args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	client, _, err := getClient(gitPullServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	gitCmd := fmt.Sprintf("cd %s && git pull", dir)
	output, err := client.Run(gitCmd)
	if err != nil {
		return fmt.Errorf("error al hacer pull: %v\n%s", err, output)
	}

	fmt.Print(output)
	return nil
}
