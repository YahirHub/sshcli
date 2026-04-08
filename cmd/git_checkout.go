package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var gitCheckoutServer string

var gitCheckoutCmd = &cobra.Command{
	Use:   "git-checkout [directorio] [rama_o_commit]",
	Short: "Cambia de rama o restaura archivos",
	Long: `Cambia a otra rama o commit en el repositorio Git remoto.

Ejemplos:
  sshcli git-checkout /app/proyecto main
  sshcli git-checkout /app/proyecto develop
  sshcli git-checkout /app/proyecto abc123f         # Checkout a commit
  sshcli git-checkout --server prod /var/www/app release-v2

Casos de uso para agentes:
  - Cambiar entre ramas
  - Restaurar versión anterior
  - Probar código de otra rama`,
	Args: cobra.ExactArgs(2),
	RunE: runGitCheckout,
}

func init() {
	rootCmd.AddCommand(gitCheckoutCmd)
	gitCheckoutCmd.Flags().StringVarP(&gitCheckoutServer, "server", "s", "", "Servidor específico a usar")
}

func runGitCheckout(cmd *cobra.Command, args []string) error {
	dir := args[0]
	target := args[1]

	client, _, err := getClient(gitCheckoutServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	gitCmd := fmt.Sprintf("cd %s && git checkout %s", dir, target)
	output, err := client.Run(gitCmd)
	if err != nil {
		return fmt.Errorf("error al hacer checkout: %v\n%s", err, output)
	}

	fmt.Print(output)
	return nil
}
