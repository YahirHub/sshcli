package cmd

import (
	"fmt"
	"path/filepath"
	"sshcli/internal/paths"
	"strings"

	"github.com/spf13/cobra"
)

var syntaxCheckServer string

var syntaxCheckCmd = &cobra.Command{
	Use:   "syntax-check [archivo]",
	Short: "Verifica la sintaxis de un archivo de código remoto",
	Args:  cobra.ExactArgs(1),
	RunE:  runSyntaxCheck,
}

func init() {
	rootCmd.AddCommand(syntaxCheckCmd)
	syntaxCheckCmd.Flags().StringVarP(&syntaxCheckServer, "server", "s", "", "Servidor específico a usar")
}

func runSyntaxCheck(cmd *cobra.Command, args []string) error {
	remotePath := paths.ToRemote(args[0])

	client, _, err := getClient(syntaxCheckServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	ext := strings.ToLower(filepath.Ext(remotePath))
	var checkCmd string

	switch ext {
	case ".py":
		checkCmd = fmt.Sprintf("python3 -m py_compile %s 2>&1", remotePath)
	case ".js":
		checkCmd = fmt.Sprintf("node --check %s 2>&1", remotePath)
	case ".ts":
		checkCmd = fmt.Sprintf("tsc --noEmit %s 2>&1 || node --check %s 2>&1", remotePath, remotePath)
	case ".go":
		checkCmd = fmt.Sprintf("cd $(dirname %s) && go build -n ./... 2>&1", remotePath)
	case ".sh", ".bash":
		checkCmd = fmt.Sprintf("bash -n %s 2>&1", remotePath)
	case ".rb":
		checkCmd = fmt.Sprintf("ruby -c %s 2>&1", remotePath)
	case ".php":
		checkCmd = fmt.Sprintf("php -l %s 2>&1", remotePath)
	case ".json":
		checkCmd = fmt.Sprintf("python3 -m json.tool %s > /dev/null 2>&1 && echo 'JSON válido' || echo 'JSON inválido'", remotePath)
	case ".yaml", ".yml":
		checkCmd = fmt.Sprintf("python3 -c \"import yaml; yaml.safe_load(open('%s'))\" 2>&1 && echo 'YAML válido' || echo 'Error en YAML'", remotePath)
	default:
		return fmt.Errorf("extensión '%s' no soportada para verificación remota", ext)
	}

	output, err := client.Run(checkCmd)
	output = strings.TrimSpace(output)
	
	if err != nil {
		fmt.Printf("ERROR de sintaxis en %s:\n", remotePath)
		if output != "" {
			fmt.Println(output)
		}
		return fmt.Errorf("sintaxis inválida")
	}

	if output == "" {
		fmt.Printf("OK: %s - Sintaxis correcta\n", remotePath)
	} else {
		fmt.Printf("Resultado para %s:\n%s\n", remotePath, output)
	}

	return nil
}