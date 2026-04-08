package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var gitCloneServer string

var gitCloneCmd = &cobra.Command{
	Use:   "git-clone [url] [destino]",
	Short: "Clona un repositorio Git",
	Long: `Clona un repositorio Git en el servidor remoto.

Ejemplos:
  sshcli git-clone https://github.com/user/repo.git /app/proyecto
  sshcli git-clone git@github.com:user/repo.git /var/www/app
  sshcli git-clone --server prod https://github.com/user/repo.git /opt/app

Casos de uso para agentes:
  - Desplegar nuevo proyecto
  - Clonar dependencias
  - Inicializar ambiente de desarrollo`,
	Args: cobra.ExactArgs(2),
	RunE: runGitClone,
}

func init() {
	rootCmd.AddCommand(gitCloneCmd)
	gitCloneCmd.Flags().StringVarP(&gitCloneServer, "server", "s", "", "Servidor específico a usar")
}

func runGitClone(cmd *cobra.Command, args []string) error {
	repoURL := args[0]
	destPath := args[1]

	client, _, err := getClient(gitCloneServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	gitCmd := fmt.Sprintf("git clone %s %s", repoURL, destPath)
	output, err := client.Run(gitCmd)
	if err != nil {
		return fmt.Errorf("error al clonar: %v\n%s", err, output)
	}

	fmt.Print(output)
	fmt.Printf("Repositorio clonado en: %s\n", destPath)
	return nil
}
