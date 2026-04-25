package cmd

import (
	"fmt"
	"strings"
	"sshcli/internal/paths"

	"github.com/spf13/cobra"
)

var (
	searchReplaceAll    bool
	searchReplaceServer string
)

var searchReplaceCmd = &cobra.Command{
	Use:   "search-replace[archivo] [buscar] [reemplazar]",
	Short: "Busca y reemplaza texto en un archivo remoto",
	Long: `Busca y reemplaza una cadena de texto dentro de un archivo en el servidor remoto.
Herramienta quirúrgica para modificar configuraciones o código.

Características:
  - Normalización automática de rutas (evita errores de Windows/Git Bash).
  - Por defecto reemplaza solo la primera ocurrencia.
  - Usa la flag --all (-a) para reemplazar todas las coincidencias.
  - Operación segura: lee el archivo, modifica en memoria y sobreescribe.

Ejemplos:
  sshcli search-replace /etc/nginx/nginx.conf "80" "8080"
  sshcli search-replace --all /app/main.py "print(" "logger.info("
  sshcli search-replace --server prod /var/www/.env "DEBUG=true" "DEBUG=false"`,
	Args: cobra.ExactArgs(3),
	RunE: runSearchReplace,
}

func init() {
	rootCmd.AddCommand(searchReplaceCmd)
	searchReplaceCmd.Flags().BoolVarP(&searchReplaceAll, "all", "a", false, "Reemplazar todas las ocurrencias encontradas")
	searchReplaceCmd.Flags().StringVarP(&searchReplaceServer, "server", "s", "", "Servidor específico a usar")
}

func runSearchReplace(cmd *cobra.Command, args[]string) error {
	remotePath := paths.ToRemote(args[0])
	search := args[1]
	replace := args[2]

	client, _, err := getClient(searchReplaceServer)
	if err != nil {
		return fmt.Errorf("error de conexión: %v", err)
	}
	defer client.Close()

	data, err := client.ReadFile(remotePath)
	if err != nil {
		return fmt.Errorf("error al leer archivo remoto en %s: %v", remotePath, err)
	}

	content := string(data)
	var newContent string
	var count int

	if searchReplaceAll {
		count = strings.Count(content, search)
		if count > 0 {
			newContent = strings.ReplaceAll(content, search, replace)
		}
	} else {
		if strings.Contains(content, search) {
			count = 1
			newContent = strings.Replace(content, search, replace, 1)
		}
	}

	if count == 0 {
		fmt.Printf("Sin coincidencias para '%s' en %s\n", search, remotePath)
		return nil
	}

	if err := client.WriteFile(remotePath,[]byte(newContent), 0644); err != nil {
		return fmt.Errorf("error al escribir los cambios en el servidor: %v", err)
	}

	fmt.Printf("[OK] Reemplazado exitosamente: %d ocurrencia(s) en %s\n", count, remotePath)
	return nil
}