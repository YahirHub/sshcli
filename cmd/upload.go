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

func runUpload(cmd *cobra.Command, args []string) error {
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

	if info.IsDir() {
		return filepath.Walk(localPath, func(p string, i os.FileInfo, err error) error {
			if err != nil || i.IsDir() {
				return nil
			}
			
			rel, _ := filepath.Rel(localPath, p)
			// Destino remoto debe usar forward slashes
			remoteDest := path.Join(remotePath, strings.ReplaceAll(rel, "\\", "/"))
			
			// Asegurar subdirectorio
			_, _ = client.Run(fmt.Sprintf("mkdir -p %s", path.Dir(remoteDest)))
			
			data, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			
			return client.WriteFile(remoteDest, data, i.Mode())
		})
	}

	// Archivo único
	data, err := os.ReadFile(localPath)
	if err != nil {
		return err
	}

	_, _ = client.Run(fmt.Sprintf("mkdir -p %s", path.Dir(remotePath)))
	return client.WriteFile(remotePath, data, 0644)
}