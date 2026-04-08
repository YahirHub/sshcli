package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	searchReplaceAll    bool
	searchReplaceServer string
)

var searchReplaceCmd = &cobra.Command{
	Use:   "search-replace [archivo] [buscar] [reemplazar]",
	Short: "Busca y reemplaza texto en un archivo remoto",
	Long: `Busca y reemplaza texto dentro de un archivo remoto.
Herramienta esencial para modificar código de forma quirúrgica.
Por defecto reemplaza solo la primera ocurrencia.

Ejemplos:
  sshcli search-replace /app/config.py "localhost" "prod.db.com"
  sshcli search-replace --all /app/main.py "print(" "logger.info("
  sshcli search-replace --server prod /etc/nginx/nginx.conf "80" "8080"
  
Casos de uso para agentes:
  - Cambiar configuraciones
  - Renombrar variables/funciones
  - Actualizar imports
  - Corregir bugs puntuales`,
	Args: cobra.ExactArgs(3),
	RunE: runSearchReplace,
}

func init() {
	rootCmd.AddCommand(searchReplaceCmd)
	searchReplaceCmd.Flags().BoolVarP(&searchReplaceAll, "all", "a", false, "Reemplazar todas las ocurrencias")
	searchReplaceCmd.Flags().StringVarP(&searchReplaceServer, "server", "s", "", "Servidor específico a usar")
}

func runSearchReplace(cmd *cobra.Command, args []string) error {
	remotePath := args[0]
	search := args[1]
	replace := args[2]

	client, _, err := getClient(searchReplaceServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	// Leer archivo
	data, err := client.ReadFile(remotePath)
	if err != nil {
		return fmt.Errorf("error al leer archivo: %v", err)
	}

	content := string(data)
	var newContent string
	var count int

	if searchReplaceAll {
		count = strings.Count(content, search)
		newContent = strings.ReplaceAll(content, search, replace)
	} else {
		if strings.Contains(content, search) {
			count = 1
			newContent = strings.Replace(content, search, replace, 1)
		} else {
			count = 0
			newContent = content
		}
	}

	if count == 0 {
		fmt.Println("Sin coincidencias encontradas")
		return nil
	}

	// Escribir archivo modificado
	if err := client.WriteFile(remotePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("error al escribir archivo: %v", err)
	}

	fmt.Printf("Reemplazado: %d ocurrencia(s) en %s\n", count, remotePath)
	return nil
}
