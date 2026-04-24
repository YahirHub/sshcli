package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var gitCloneServer string

var gitCloneCmd = &cobra.Command{
	Use:   "git-clone [url_o_ruta] [destino]",
	Short: "Clona un repositorio Git",
	Args:  cobra.ExactArgs(2),
	RunE:  runGitClone,
}

func init() {
	rootCmd.AddCommand(gitCloneCmd)
	gitCloneCmd.Flags().StringVarP(&gitCloneServer, "server", "s", "", "Servidor específico a usar")
}

func runGitClone(cmd *cobra.Command, args []string) error {
	// IMPORTANTE: Limpiamos también la URL por si es una ruta local del servidor
	// cleanRemotePath es seguro incluso con URLs reales (http, git@, etc)
	repoURL := cleanRemotePath(args[0])
	destPath := cleanRemotePath(args[1])

	client, _, err := getClient(gitCloneServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	// Usamos comillas simples para proteger ambas rutas en el servidor remoto
	gitCmd := fmt.Sprintf("git clone '%s' '%s'", repoURL, destPath)
	output, err := client.Run(gitCmd)
	if err != nil {
		return fmt.Errorf("error al clonar: %v\n%s", err, output)
	}

	fmt.Print(output)
	fmt.Printf("Repositorio clonado en: %s\n", destPath)
	return nil
}