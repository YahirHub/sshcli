package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var gitAddServer string

var gitAddCmd = &cobra.Command{
	Use:   "git-add [directorio] [archivos...]",
	Short: "Agrega archivos al staging area",
	Long: `Agrega archivos al área de staging para el próximo commit.
Usa '.' para agregar todos los archivos.

Ejemplos:
  sshcli git-add /app/proyecto .                    # Agregar todo
  sshcli git-add /app/proyecto main.py config.py   # Archivos específicos
  sshcli git-add --server prod /var/www/app .

Casos de uso para agentes:
  - Preparar archivos para commit
  - Seleccionar cambios específicos
  - Stage de nuevos archivos`,
	Args: cobra.MinimumNArgs(2),
	RunE: runGitAdd,
}

func init() {
	rootCmd.AddCommand(gitAddCmd)
	gitAddCmd.Flags().StringVarP(&gitAddServer, "server", "s", "", "Servidor específico a usar")
}

func runGitAdd(cmd *cobra.Command, args []string) error {
	dir := args[0]
	files := args[1:]

	client, _, err := getClient(gitAddServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	fileList := ""
	for _, f := range files {
		fileList += f + " "
	}

	gitCmd := fmt.Sprintf("cd %s && git add %s", dir, fileList)
	output, err := client.Run(gitCmd)
	if err != nil {
		return fmt.Errorf("error al agregar archivos: %v\n%s", err, output)
	}

	if output != "" {
		fmt.Print(output)
	}
	fmt.Println("Archivos agregados al staging")
	return nil
}
