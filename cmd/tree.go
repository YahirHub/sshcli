package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	treeDepth   int
	treeDirs    bool
	treeServer  string
)

var treeCmd = &cobra.Command{
	Use:   "tree [directorio]",
	Short: "Muestra la estructura de directorios en forma de árbol",
	Long: `Muestra la estructura de archivos y directorios en formato árbol.
Esencial para entender la organización de un proyecto.

Ejemplos:
  sshcli tree /app
  sshcli tree /var/www -d 2              # Profundidad máxima 2
  sshcli tree /app --dirs                 # Solo directorios
  sshcli tree --server prod /opt/app

Casos de uso para agentes:
  - Entender estructura de proyecto
  - Localizar archivos de configuración
  - Mapear arquitectura de código
  - Identificar módulos y componentes`,
	Args: cobra.MaximumNArgs(1),
	RunE: runTree,
}

func init() {
	rootCmd.AddCommand(treeCmd)
	treeCmd.Flags().IntVarP(&treeDepth, "depth", "d", 3, "Profundidad máxima del árbol")
	treeCmd.Flags().BoolVar(&treeDirs, "dirs", false, "Mostrar solo directorios")
	treeCmd.Flags().StringVarP(&treeServer, "server", "s", "", "Servidor específico a usar")
}

func runTree(cmd *cobra.Command, args []string) error {
	remotePath := "."
	if len(args) > 0 {
		remotePath = args[0]
	}

	client, _, err := getClient(treeServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	// Primero intentamos con tree si está instalado
	var treeCommand string
	if treeDirs {
		treeCommand = fmt.Sprintf("tree -d -L %d %s 2>/dev/null || find %s -maxdepth %d -type d | head -100", 
			treeDepth, remotePath, remotePath, treeDepth)
	} else {
		treeCommand = fmt.Sprintf("tree -L %d %s 2>/dev/null || find %s -maxdepth %d | head -200", 
			treeDepth, remotePath, remotePath, treeDepth)
	}

	output, err := client.Run(treeCommand)
	if err != nil {
		// Fallback a find con formato
		fallbackCmd := fmt.Sprintf("find %s -maxdepth %d 2>/dev/null | sort | head -200", remotePath, treeDepth)
		output, err = client.Run(fallbackCmd)
		if err != nil {
			return fmt.Errorf("error al listar estructura: %v", err)
		}
	}

	if output == "" {
		fmt.Printf("Directorio vacío o no accesible: %s\n", remotePath)
		return nil
	}

	fmt.Print(output)
	return nil
}
