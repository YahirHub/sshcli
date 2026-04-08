package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"sshcli/internal/ssh"

	"github.com/spf13/cobra"
)

var downloadServer string

var downloadCmd = &cobra.Command{
	Use:   "download [origen_remoto] [destino_local]",
	Short: "Descarga un archivo o carpeta del servidor remoto",
	Long: `Descarga un archivo o carpeta completa del servidor remoto.
Si el origen es una carpeta, se descarga recursivamente.

Ejemplos:
  sshcli download /home/user/archivo.txt ./archivo.txt
  sshcli download --server produccion /var/log/app.log ./logs/app.log
  sshcli download /home/user/proyecto ./proyecto_local`,
	Args: cobra.ExactArgs(2),
	RunE: runDownload,
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVarP(&downloadServer, "server", "s", "", "Servidor específico a usar")
}

func runDownload(cmd *cobra.Command, args []string) error {
	remotePath := args[0]
	localPath := args[1]

	client, _, err := getClient(downloadServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	output, err := client.Run(fmt.Sprintf("test -d %s && echo 'dir' || echo 'file'", remotePath))
	if err != nil {
		return fmt.Errorf("error al verificar ruta remota: %v", err)
	}

	isDir := strings.TrimSpace(output) == "dir"

	if isDir {
		return downloadDirectoryNew(client, remotePath, localPath)
	}

	return downloadFileNew(client, remotePath, localPath)
}

func downloadFileNew(client *ssh.Client, remotePath, localPath string) error {
	data, err := client.ReadFile(remotePath)
	if err != nil {
		return fmt.Errorf("error al leer archivo remoto: %v", err)
	}

	dir := filepath.Dir(localPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error al crear directorio local: %v", err)
	}

	if err := os.WriteFile(localPath, data, 0644); err != nil {
		return fmt.Errorf("error al escribir archivo local: %v", err)
	}

	fmt.Printf("Archivo descargado: %s -> %s\n", remotePath, localPath)
	return nil
}

func downloadDirectoryNew(client *ssh.Client, remoteDir, localDir string) error {
	if err := os.MkdirAll(localDir, 0755); err != nil {
		return fmt.Errorf("error al crear directorio local: %v", err)
	}

	output, err := client.Run(fmt.Sprintf("find %s -type f", remoteDir))
	if err != nil {
		return fmt.Errorf("error al listar archivos remotos: %v", err)
	}

	files := strings.Split(strings.TrimSpace(output), "\n")
	for _, file := range files {
		if file == "" {
			continue
		}

		relPath, _ := filepath.Rel(remoteDir, file)
		localDest := filepath.Join(localDir, relPath)

		data, err := client.ReadFile(file)
		if err != nil {
			fmt.Printf("Error al descargar %s: %v\n", file, err)
			continue
		}

		dir := filepath.Dir(localDest)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		if err := os.WriteFile(localDest, data, 0644); err != nil {
			return err
		}

		fmt.Printf("Descargado: %s -> %s\n", file, localDest)
	}

	fmt.Printf("Carpeta descargada exitosamente: %s -> %s\n", remoteDir, localDir)
	return nil
}
