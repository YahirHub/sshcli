package cmd

import (
	"fmt"
	"sshcli/internal/paths"

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
	Args:  cobra.MaximumNArgs(1),
	RunE:  runTree,
}

func init() {
	rootCmd.AddCommand(treeCmd)
	treeCmd.Flags().IntVarP(&treeDepth, "depth", "d", 3, "Profundidad máxima")
	treeCmd.Flags().BoolVar(&treeDirs, "dirs", false, "Solo directorios")
	treeCmd.Flags().StringVarP(&treeServer, "server", "s", "", "Servidor específico a usar")
}

func runTree(cmd *cobra.Command, args[]string) error {
	remotePath := "/"
	if len(args) > 0 {
		remotePath = paths.ToRemote(args[0])
	}

	client, _, err := getClient(treeServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	var treeCommand string
	if treeDirs {
		treeCommand = fmt.Sprintf("tree -d -L %d '%s' 2>/dev/null || find '%s' -maxdepth %d -type d | head -100", 
			treeDepth, remotePath, remotePath, treeDepth)
	} else {
		treeCommand = fmt.Sprintf("tree -L %d '%s' 2>/dev/null || find '%s' -maxdepth %d | head -200", 
			treeDepth, remotePath, remotePath, treeDepth)
	}

	output, err := client.Run(treeCommand)
	if err != nil {
		// Fallback manual si 'tree' no está instalado
		fallbackCmd := fmt.Sprintf("find '%s' -maxdepth %d 2>/dev/null | sort | head -200", remotePath, treeDepth)
		output, _ = client.Run(fallbackCmd)
	}

	if output == "" {
		fmt.Printf("Estructura de %s (vía find):\n", remotePath)
	}
	fmt.Print(output)
	return nil
}