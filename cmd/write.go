package cmd

import (
	"fmt"
	"io"
	"os"
	"path"
	"sshcli/internal/paths"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	writeServer string
	writeChmod  string
	writeExec   bool
)

var writeCmd = &cobra.Command{
	Use:   "write [ruta_remota] [contenido_opcional]",
	Short: "Escribe contenido a un archivo remoto",
	Long: `Escribe contenido directamente a un archivo remoto.
Garantiza la creación de directorios y permite permisos -x para scripts.`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runWrite,
}

func init() {
	rootCmd.AddCommand(writeCmd)
	writeCmd.Flags().StringVarP(&writeServer, "server", "s", "", "Servidor específico a usar")
	writeCmd.Flags().StringVar(&writeChmod, "chmod", "644", "Permisos octales (ej: 644)")
	writeCmd.Flags().BoolVarP(&writeExec, "exec", "x", false, "Hacer ejecutable (755)")
}

func runWrite(cmd *cobra.Command, args[]string) error {
	remotePath := paths.ToRemote(args[0])
	
	var content[]byte
	var err error

	if len(args) == 2 {
		content =[]byte(args[1])
	} else {
		content, err = io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("error al leer stdin: %v", err)
		}
	}

	client, _, err := getClient(writeServer)
	if err != nil {
		return err
	}
	defer client.Close()

	// Procesar permisos
	permStr := writeChmod
	if writeExec {
		permStr = "755"
	}
	p, err := strconv.ParseUint(permStr, 8, 32)
	if err != nil {
		return fmt.Errorf("error en permisos '%s': %v", permStr, err)
	}

	// Asegurar carpeta remota
	dir := path.Dir(remotePath)
	if dir != "." && dir != "/" {
		if _, err := client.Run(fmt.Sprintf("mkdir -p '%s'", dir)); err != nil {
			return fmt.Errorf("error al crear directorio: %v", err)
		}
	}

	if err := client.WriteFile(remotePath, content, os.FileMode(p)); err != nil {
		return fmt.Errorf("error al escribir: %v", err)
	}

	fmt.Printf("[OK] Archivo escrito exitosamente: %s (mode: %s, size: %d bytes)\n", remotePath, permStr, len(content))
	return nil
}