package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var uploadServer string

var uploadCmd = &cobra.Command{
	Use:   "upload [origen_local] [destino_remoto]",
	Short: "Sube un archivo o carpeta al servidor remoto",
	Long: `Sube un archivo o carpeta completa al servidor remoto.
Si el origen es una carpeta, se sube recursivamente.

Ejemplos:
  sshcli upload archivo.txt /home/user/archivo.txt
  sshcli upload --server produccion ./proyecto /home/user/proyecto
  sshcli upload config.json /etc/app/config.json`,
	Args: cobra.ExactArgs(2),
	RunE: runUpload,
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().StringVarP(&uploadServer, "server", "s", "", "Servidor específico a usar")
}

func runUpload(cmd *cobra.Command, args []string) error {
	localPath := args[0]
	remotePath := args[1]

	client, _, err := getClient(uploadServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	info, err := os.Stat(localPath)
	if err != nil {
		return fmt.Errorf("error al acceder a '%s': %v", localPath, err)
	}

	if info.IsDir() {
		return uploadDirectoryNew(client, localPath, remotePath)
	}

	return uploadFileNew(client, localPath, remotePath)
}

func uploadFileNew(client interface{ WriteFile(string, []byte, os.FileMode) error; Run(string) (string, error) }, localPath, remotePath string) error {
	data, err := os.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("error al leer archivo: %v", err)
	}

	if err := client.WriteFile(remotePath, data, 0644); err != nil {
		return fmt.Errorf("error al subir archivo: %v", err)
	}

	fmt.Printf("Archivo subido: %s -> %s\n", localPath, remotePath)
	return nil
}

func uploadDirectoryNew(client interface{ WriteFile(string, []byte, os.FileMode) error; Run(string) (string, error) }, localDir, remoteDir string) error {
	if _, err := client.Run(fmt.Sprintf("mkdir -p %s", remoteDir)); err != nil {
		return fmt.Errorf("error al crear directorio remoto: %v", err)
	}

	err := filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(localDir, path)
		remoteDest := filepath.Join(remoteDir, relPath)

		if info.IsDir() {
			_, err := client.Run(fmt.Sprintf("mkdir -p %s", remoteDest))
			return err
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if err := client.WriteFile(remoteDest, data, info.Mode()); err != nil {
			return err
		}

		fmt.Printf("Subido: %s -> %s\n", path, remoteDest)
		return nil
	})

	if err != nil {
		return fmt.Errorf("error al subir carpeta: %v", err)
	}

	fmt.Printf("Carpeta subida exitosamente: %s -> %s\n", localDir, remoteDir)
	return nil
}
