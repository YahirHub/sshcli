package cmd

import (
	"fmt"
	"os"
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
	Long: `Descarga un archivo o carpeta completa del servidor remoto.
Es inteligente: detecta si el destino es un directorio y coloca el archivo dentro.

Características:
  - Descarga archivos individuales o carpetas completas
  - Detecta automáticamente si el destino local es una carpeta
  - Sincronización con exclusiones personalizables
  - Modo dry-run para previsualizar cambios

Ejemplos:
  sshcli download /home/user/archivo.txt ./              # Descarga a ./archivo.txt
  sshcli download /tmp/backup.zip ./backups/             # Descarga a ./backups/backup.zip
  sshcli download /var/www/app ./app_local               # Descarga carpeta completa
  sshcli download /srv/project ./local --sync            # Sincroniza carpeta
  sshcli download /srv/app ./app --exclude ".git" --exclude "node_modules"
  sshcli download /srv/app ./app --exclude ".git,*.log,tmp/"
  sshcli download /tmp/data.tar ./backups/ --dry-run     # Solo muestra qué haría`,
	Args: cobra.ExactArgs(2),
	RunE: runDownload,
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVarP(&downloadServer, "server", "s", "", "Servidor específico a usar")
	downloadCmd.Flags().StringArrayVarP(&downloadExclude, "exclude", "e", []string{}, "Patrones a excluir (puede usarse múltiples veces o separados por coma)")
	downloadCmd.Flags().BoolVar(&downloadSync, "sync", false, "Modo sincronización (elimina archivos locales que no existen en remoto)")
	downloadCmd.Flags().BoolVar(&downloadDryRun, "dry-run", false, "Solo muestra qué archivos se descargarían sin ejecutar")
}

func runDownload(cmd *cobra.Command, args []string) error {
	remotePath := args[0]
	localPath := args[1]

	client, _, err := getClient(downloadServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	// Procesar exclusiones
	excludePatterns := processDownloadExcludes(downloadExclude)

	// Verificar si el remoto es directorio
	output, err := client.Run(fmt.Sprintf("test -d %s && echo 'dir' || echo 'file'", remotePath))
	if err != nil {
		return fmt.Errorf("error al verificar ruta remota: %v", err)
	}
	isRemoteDir := strings.TrimSpace(output) == "dir"

	// Detectar si el destino local es un directorio
	isLocalDir := false
	if strings.HasSuffix(localPath, "/") || strings.HasSuffix(localPath, string(os.PathSeparator)) {
		isLocalDir = true
	} else if info, err := os.Stat(localPath); err == nil && info.IsDir() {
		isLocalDir = true
	}

	// Si es archivo remoto y destino local es directorio, agregar nombre del archivo
	if !isRemoteDir && isLocalDir {
		localPath = filepath.Join(localPath, filepath.Base(remotePath))
	}

	if downloadDryRun {
		fmt.Println("=== MODO DRY-RUN (sin cambios reales) ===")
	}

	if isRemoteDir {
		return downloadDirectoryNew(client, remotePath, localPath, excludePatterns, downloadSync, downloadDryRun)
	}

	return downloadFileNewSmart(client, remotePath, localPath, downloadDryRun)
}

// processDownloadExcludes procesa los patrones de exclusión
func processDownloadExcludes(excludes []string) []string {
	var patterns []string
	for _, exc := range excludes {
		parts := strings.Split(exc, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				patterns = append(patterns, p)
			}
		}
	}
	return patterns
}

// shouldExcludeDownload verifica si una ruta debe ser excluida
func shouldExcludeDownload(path string, patterns []string) bool {
	baseName := filepath.Base(path)
	for _, pattern := range patterns {
		if baseName == pattern || baseName == strings.TrimSuffix(pattern, "/") {
			return true
		}
		if matched, _ := filepath.Match(pattern, baseName); matched {
			return true
		}
		if strings.Contains(path, pattern) {
			return true
		}
	}
	return false
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
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error al crear directorio local: %v", err)
	}

	if err := os.WriteFile(localPath, data, 0644); err != nil {
		return fmt.Errorf("error al escribir archivo local: %v", err)
	}

	fmt.Printf("✓ Descargado: %s -> %s\n", remotePath, localPath)
	return nil
}

func downloadDirectoryNew(client *ssh.Client, remoteDir, localDir string, excludePatterns []string, sync bool, dryRun bool) error {
	if !dryRun {
		if err := os.MkdirAll(localDir, 0755); err != nil {
			return fmt.Errorf("error al crear directorio local: %v", err)
		}
	} else {
		fmt.Printf("[DRY-RUN] Crearía directorio: %s\n", localDir)
	}

	output, err := client.Run(fmt.Sprintf("find %s -type f 2>/dev/null", remoteDir))
	if err != nil {
		return fmt.Errorf("error al listar archivos remotos: %v", err)
	}

	downloadedFiles := make(map[string]bool)
	var filesCount, skippedCount int

	files := strings.Split(strings.TrimSpace(output), "\n")
	for _, file := range files {
		if file == "" {
			continue
		}

		relPath, _ := filepath.Rel(remoteDir, file)

		// Verificar exclusiones
		if shouldExcludeDownload(relPath, excludePatterns) || shouldExcludeDownload(file, excludePatterns) {
			fmt.Printf("⊘ Omitido (excluido): %s\n", relPath)
			skippedCount++
			continue
		}

		localDest := filepath.Join(localDir, relPath)
		downloadedFiles[relPath] = true

		if dryRun {
			fmt.Printf("[DRY-RUN] Descargaría: %s -> %s\n", file, localDest)
			filesCount++
			continue
		}

		data, err := client.ReadFile(file)
		if err != nil {
			fmt.Printf("✗ Error al descargar %s: %v\n", relPath, err)
			continue
		}

		dir := filepath.Dir(localDest)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		if err := os.WriteFile(localDest, data, 0644); err != nil {
			return err
		}

		fmt.Printf("✓ Descargado: %s\n", relPath)
		filesCount++
	}

	// Sincronización: eliminar archivos locales que no existen en remoto
	if sync && !dryRun {
		fmt.Println("\n--- Sincronizando (eliminando archivos locales extras) ---")
		err := filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			relPath, _ := filepath.Rel(localDir, path)
			if !downloadedFiles[relPath] && !shouldExcludeDownload(relPath, excludePatterns) {
				os.Remove(path)
				fmt.Printf("✗ Eliminado (sync): %s\n", relPath)
			}
			return nil
		})
		if err != nil {
			fmt.Printf("Advertencia en sincronización: %v\n", err)
		}
	}

	fmt.Printf("\n=== Resumen ===\n")
	fmt.Printf("Carpeta: %s -> %s\n", remoteDir, localDir)
	fmt.Printf("Archivos: %d, Omitidos: %d\n", filesCount, skippedCount)
	if dryRun {
		fmt.Println("(Modo dry-run: no se realizaron cambios)")
	}

	return nil
}
