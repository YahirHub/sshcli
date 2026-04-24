package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"sshcli/internal/ssh"

	"github.com/spf13/cobra"
)

var (
	downloadServer  string
	downloadExclude []string
	downloadSync    bool
	downloadDryRun  bool
)

var downloadCmd = &cobra.Command{
	Use:   "download [origen_remoto] [destino_local]",
	Short: "Descarga un archivo o carpeta del servidor remoto",
	Long: `Descarga archivos. Soporta rutas locales estilo /c/... para Windows.
Limpia automáticamente rutas remotas si el shell las mutila.`,
	Args: cobra.ExactArgs(2),
	RunE: runDownload,
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVarP(&downloadServer, "server", "s", "", "Servidor específico a usar")
	downloadCmd.Flags().StringArrayVarP(&downloadExclude, "exclude", "e", []string{}, "Patrones a excluir")
	downloadCmd.Flags().BoolVar(&downloadSync, "sync", false, "Modo sincronización")
	downloadCmd.Flags().BoolVar(&downloadDryRun, "dry-run", false, "Modo dry-run")
}

func runDownload(cmd *cobra.Command, args []string) error {
	remotePath := cleanRemotePath(args[0])
	localPath := cleanLocalPath(args[1])

	client, _, err := getClient(downloadServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	output, err := client.Run(fmt.Sprintf("test -d %s && echo 'dir' || echo 'file'", remotePath))
	if err != nil {
		return fmt.Errorf("error al verificar ruta remota: %v", err)
	}
	isRemoteDir := strings.TrimSpace(output) == "dir"

	// Si el local es un directorio (termina en slash o existe)
	isLocalDir := strings.HasSuffix(localPath, string(os.PathSeparator))
	if !isLocalDir {
		if info, err := os.Stat(localPath); err == nil && info.IsDir() {
			isLocalDir = true
		}
	}

	if !isRemoteDir && isLocalDir {
		localPath = filepath.Join(localPath, path.Base(remotePath))
	}

	if downloadDryRun {
		fmt.Println("=== MODO DRY-RUN ===")
	}

	if isRemoteDir {
		return downloadDirectoryNew(client, remotePath, localPath, processDownloadExcludes(downloadExclude), downloadSync, downloadDryRun)
	}

	return downloadFileNewSmart(client, remotePath, localPath, downloadDryRun)
}

func downloadFileNewSmart(client *ssh.Client, remotePath, localPath string, dryRun bool) error {
	if dryRun {
		fmt.Printf("[DRY-RUN] Descargaría: %s -> %s\n", remotePath, localPath)
		return nil
	}

	data, err := client.ReadFile(remotePath)
	if err != nil {
		return fmt.Errorf("error al leer archivo remoto: %v", err)
	}

	dir := filepath.Dir(localPath)
	os.MkdirAll(dir, 0755)

	if err := os.WriteFile(localPath, data, 0644); err != nil {
		return fmt.Errorf("error al escribir archivo local: %v", err)
	}

	fmt.Printf("✓ Descargado: %s -> %s\n", remotePath, localPath)
	return nil
}

func downloadDirectoryNew(client *ssh.Client, remoteDir, localDir string, excludePatterns []string, sync bool, dryRun bool) error {
	if !dryRun {
		os.MkdirAll(localDir, 0755)
	}

	output, err := client.Run(fmt.Sprintf("find %s -type f 2>/dev/null", remoteDir))
	if err != nil {
		return err
	}

	files := strings.Split(strings.TrimSpace(output), "\n")
	for _, f := range files {
		if f == "" {
			continue
		}
		rel, _ := filepath.Rel(remoteDir, f)
		relLocal := strings.ReplaceAll(rel, "/", string(os.PathSeparator))
		localDest := filepath.Join(localDir, relLocal)

		if dryRun {
			fmt.Printf("[DRY-RUN] %s -> %s\n", f, localDest)
			continue
		}

		data, err := client.ReadFile(f)
		if err != nil {
			continue
		}

		os.MkdirAll(filepath.Dir(localDest), 0755)
		os.WriteFile(localDest, data, 0644)
		fmt.Printf("✓ Descargado: %s\n", rel)
	}
	return nil
}

func processDownloadExcludes(excludes []string) []string {
	var patterns []string
	for _, exc := range excludes {
		parts := strings.Split(exc, ",")
		for _, p := range parts {
			patterns = append(patterns, strings.TrimSpace(p))
		}
	}
	return patterns
}