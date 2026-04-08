package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var syntaxCheckServer string

var syntaxCheckCmd = &cobra.Command{
	Use:   "syntax-check [archivo]",
	Short: "Verifica la sintaxis de un archivo de código",
	Long: `Verifica la sintaxis de archivos de código sin ejecutarlos.
Soporta Python, JavaScript/Node, Go, Bash, Ruby, PHP.
Detecta el lenguaje automáticamente por la extensión.

Ejemplos:
  sshcli syntax-check /app/main.py
  sshcli syntax-check /app/index.js
  sshcli syntax-check --server dev /app/main.go
  sshcli syntax-check /opt/scripts/deploy.sh

Lenguajes soportados:
  .py     -> python -m py_compile
  .js     -> node --check
  .ts     -> tsc --noEmit (si disponible)
  .go     -> go build -n
  .sh     -> bash -n
  .rb     -> ruby -c
  .php    -> php -l
  .json   -> python -m json.tool
  .yaml   -> python -c "import yaml"`,
	Args: cobra.ExactArgs(1),
	RunE: runSyntaxCheck,
}

func init() {
	rootCmd.AddCommand(syntaxCheckCmd)
	syntaxCheckCmd.Flags().StringVarP(&syntaxCheckServer, "server", "s", "", "Servidor específico a usar")
}

func runSyntaxCheck(cmd *cobra.Command, args []string) error {
	remotePath := args[0]

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
		checkCmd = fmt.Sprintf("cd $(dirname %s) && go build -n ./... 2>&1 | head -20", remotePath)
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
	case ".xml":
		checkCmd = fmt.Sprintf("python3 -c \"import xml.etree.ElementTree as ET; ET.parse('%s')\" 2>&1 && echo 'XML válido' || echo 'Error en XML'", remotePath)
	default:
		return fmt.Errorf("extensión '%s' no soportada para verificación de sintaxis", ext)
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
