package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sshcli/internal/paths"
	"strings"

	"github.com/spf13/cobra"
)

var (
	uploadServer string
	uploadSync   bool
)

var uploadCmd = &cobra.Command{
	Use:   "upload [local] [remoto]",
	Short: "Sube un archivo o carpeta al servidor remoto",
	Args:  cobra.ExactArgs(2),
	RunE:  runUpload,
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().StringVarP(&uploadServer, "server", "s", "", "Servidor específico a usar")
	uploadCmd.Flags().BoolVar(&uploadSync, "sync", false, "Modo sincronización")
}

func runUpload(cmd *cobra.Command, args[]string) error {
	localPath := paths.ToLocal(args[0])
	remotePath := paths.ToRemote(args[1])

	client, _, err := getClient(uploadServer)
	if err != nil {
		return err
	}
	defer client.Close()

	info, err := os.Stat(localPath)
	if err != nil {
		return fmt.Errorf("error en local: %v", err)
	}

	// Verificar si el destino remoto debe ser tratado como directorio
	isRemoteDir := strings.HasSuffix(args[1], "/") || strings.HasSuffix(args[1], "\\")
	if !isRemoteDir {
		// Envolver remotePath en comillas por seguridad
		output, _ := client.Run(fmt.Sprintf("test -d '%s' && echo 'dir'", remotePath))
		isRemoteDir = strings.TrimSpace(output) == "dir"
	}

	if info.IsDir() {
		// Si es directorio local y destino es directorio, la raíz será remotePath
		return filepath.Walk(localPath, func(p string, i os.FileInfo, err error) error {
			if err != nil || i.IsDir() {
				return nil
			}
			
			rel, _ := filepath.Rel(localPath, p)
			// Destino remoto debe usar forward slashes
			remoteDest := path.Join(remotePath, strings.ReplaceAll(rel, "\\", "/"))
			
			// Asegurar subdirectorio
			if _, err := client.Run(fmt.Sprintf("mkdir -p '%s'", path.Dir(remoteDest))); err != nil {
				return fmt.Errorf("error al crear subdirectorio: %v", err)
			}
			
			data, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			
			if err := client.WriteFile(remoteDest, data, i.Mode()); err == nil {
				fmt.Printf("[OK] Subido: %s -> %s\n", p, remoteDest)
			}
			return nil
		})
	}

	// Archivo único
	if isRemoteDir {
		remotePath = path.Join(remotePath, filepath.Base(localPath))
	}

	data, err := os.ReadFile(localPath)
	if err != nil {
		return err
	}

	if _, err := client.Run(fmt.Sprintf("mkdir -p '%s'", path.Dir(remotePath))); err != nil {
		return fmt.Errorf("error al crear directorio destino: %v", err)
	}
	
	if err := client.WriteFile(remotePath, data, 0644); err != nil {
		return fmt.Errorf("error al escribir archivo: %v", err)
	}

	fmt.Printf("[OK] Archivo subido exitosamente: %s -> %s\n", localPath, remotePath)
	return nil
}