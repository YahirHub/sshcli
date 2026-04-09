package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	uploadServer  string
	uploadExclude []string
	uploadSync    bool
	uploadDryRun  bool
)

var uploadCmd = &cobra.Command{
	Use:   "upload [origen_local] [destino_remoto]",
	Short: "Sube un archivo o carpeta al servidor remoto",
	Long: `Sube un archivo o carpeta completa al servidor remoto.
Es inteligente: detecta si el destino es un directorio y coloca el archivo dentro.

Características:
  - Sube archivos individuales o carpetas completas
  - Detecta automáticamente si el destino es una carpeta
  - Sincronización con exclusiones personalizables
  - Modo dry-run para previsualizar cambios

Ejemplos:
  sshcli upload archivo.txt /tmp/                    # Sube a /tmp/archivo.txt
  sshcli upload archivo.zip /tmp/                    # Sube a /tmp/archivo.zip
  sshcli upload ./proyecto /home/user/proyecto       # Sube carpeta completa
  sshcli upload ./src /home/user/ --sync             # Sincroniza carpeta
  sshcli upload ./app /srv/ --exclude ".git" --exclude "node_modules"
  sshcli upload ./app /srv/ --exclude ".git,*.log,tmp/"
  sshcli upload archivo.txt /tmp/ --dry-run          # Solo muestra qué haría`,
	Args: cobra.ExactArgs(2),
	RunE: runUpload,
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().StringVarP(&uploadServer, "server", "s", "", "Servidor específico a usar")
	uploadCmd.Flags().StringArrayVarP(&uploadExclude, "exclude", "e", []string{}, "Patrones a excluir (puede usarse múltiples veces o separados por coma)")
	uploadCmd.Flags().BoolVar(&uploadSync, "sync", false, "Modo sincronización (elimina archivos remotos que no existen localmente)")
	uploadCmd.Flags().BoolVar(&uploadDryRun, "dry-run", false, "Solo muestra qué archivos se subirían sin ejecutar")
}

func runUpload(cmd *cobra.Command, args []string) error {
	localPath := args[0]
	remotePath := args[1]

	client, _, err := getClient(uploadServer)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer client.Close()

	// Procesar exclusiones (soporta tanto --exclude ".git" --exclude "node_modules" como --exclude ".git,node_modules")
	excludePatterns := processExcludes(uploadExclude)

	localInfo, err := os.Stat(localPath)
	if err != nil {
		return fmt.Errorf("error al acceder a '%s': %v", localPath, err)
	}

	// Detectar si el destino remoto es un directorio
	isRemoteDir := false
	if strings.HasSuffix(remotePath, "/") {
		isRemoteDir = true
	} else {
		// Verificar si existe y es directorio
		output, err := client.Run(fmt.Sprintf("test -d %s && echo 'DIR' || echo 'NOTDIR'", remotePath))
		if err == nil && strings.TrimSpace(output) == "DIR" {
			isRemoteDir = true
		}
	}

	// Si es archivo y destino es directorio, agregar nombre del archivo
	if !localInfo.IsDir() && isRemoteDir {
		remotePath = filepath.Join(remotePath, filepath.Base(localPath))
	}

	if uploadDryRun {
		fmt.Println("=== MODO DRY-RUN (sin cambios reales) ===")
	}

	if localInfo.IsDir() {
		return uploadDirectoryNew(client, localPath, remotePath, excludePatterns, uploadSync, uploadDryRun)
	}

	return uploadFileNew(client, localPath, remotePath, uploadDryRun)
}

// processExcludes procesa los patrones de exclusión
func processExcludes(excludes []string) []string {
	var patterns []string
	for _, exc := range excludes {
		// Soportar separación por comas
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

// shouldExclude verifica si una ruta debe ser excluida
func shouldExclude(path string, patterns []string) bool {
	baseName := filepath.Base(path)
	for _, pattern := range patterns {
		// Coincidencia exacta con nombre de archivo/directorio
		if baseName == pattern || baseName == strings.TrimSuffix(pattern, "/") {
			return true
		}
		// Coincidencia con patrón glob
		if matched, _ := filepath.Match(pattern, baseName); matched {
			return true
		}
		// Coincidencia si el path contiene el patrón
		if strings.Contains(path, pattern) {
			return true
		}
	}
	return false
}

func uploadFileNew(client interface{ WriteFile(string, []byte, os.FileMode) error; Run(string) (string, error) }, localPath, remotePath string, dryRun bool) error {
	if dryRun {
		fmt.Printf("[DRY-RUN] Subiría: %s -> %s\n", localPath, remotePath)
		return nil
	}

	data, err := os.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("error al leer archivo: %v", err)
	}

	// Asegurar que el directorio padre existe
	parentDir := filepath.Dir(remotePath)
	if _, err := client.Run(fmt.Sprintf("mkdir -p %s", parentDir)); err != nil {
		return fmt.Errorf("error al crear directorio padre: %v", err)
	}

	if err := client.WriteFile(remotePath, data, 0644); err != nil {
		return fmt.Errorf("error al subir archivo: %v", err)
	}

	fmt.Printf("✓ Subido: %s -> %s\n", localPath, remotePath)
	return nil
}

func uploadDirectoryNew(client interface{ WriteFile(string, []byte, os.FileMode) error; Run(string) (string, error) }, localDir, remoteDir string, excludePatterns []string, sync bool, dryRun bool) error {
	if !dryRun {
		if _, err := client.Run(fmt.Sprintf("mkdir -p %s", remoteDir)); err != nil {
			return fmt.Errorf("error al crear directorio remoto: %v", err)
		}
	} else {
		fmt.Printf("[DRY-RUN] Crearía directorio: %s\n", remoteDir)
	}

	// Rastrear archivos subidos para sincronización
	uploadedFiles := make(map[string]bool)
	var filesCount, dirsCount, skippedCount int

	err := filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(localDir, path)
		if relPath == "." {
			return nil
		}

		// Verificar exclusiones
		if shouldExclude(relPath, excludePatterns) || shouldExclude(path, excludePatterns) {
			if info.IsDir() {
				fmt.Printf("⊘ Omitido (excluido): %s/\n", relPath)
				return filepath.SkipDir
			}
			fmt.Printf("⊘ Omitido (excluido): %s\n", relPath)
			skippedCount++
			return nil
		}

		remoteDest := filepath.Join(remoteDir, relPath)
		uploadedFiles[relPath] = true

		if info.IsDir() {
			if !dryRun {
				_, err := client.Run(fmt.Sprintf("mkdir -p %s", remoteDest))
				if err != nil {
					return err
				}
			} else {
				fmt.Printf("[DRY-RUN] Crearía directorio: %s\n", remoteDest)
			}
			dirsCount++
			return nil
		}

		if dryRun {
			fmt.Printf("[DRY-RUN] Subiría: %s -> %s\n", path, remoteDest)
			filesCount++
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if err := client.WriteFile(remoteDest, data, info.Mode()); err != nil {
			return err
		}

		fmt.Printf("✓ Subido: %s\n", relPath)
		filesCount++
		return nil
	})

	if err != nil {
		return fmt.Errorf("error al subir carpeta: %v", err)
	}

	// Sincronización: eliminar archivos remotos que no existen localmente
	if sync && !dryRun {
		fmt.Println("\n--- Sincronizando (eliminando archivos remotos extras) ---")
		output, err := client.Run(fmt.Sprintf("find %s -type f 2>/dev/null", remoteDir))
		if err == nil && output != "" {
			remoteFiles := strings.Split(strings.TrimSpace(output), "\n")
			for _, remoteFile := range remoteFiles {
				if remoteFile == "" {
					continue
				}
				relPath, _ := filepath.Rel(remoteDir, remoteFile)
				if !uploadedFiles[relPath] && !shouldExclude(relPath, excludePatterns) {
					client.Run(fmt.Sprintf("rm -f %s", remoteFile))
					fmt.Printf("✗ Eliminado (sync): %s\n", relPath)
				}
			}
		}
	}

	fmt.Printf("\n=== Resumen ===\n")
	fmt.Printf("Carpeta: %s -> %s\n", localDir, remoteDir)
	fmt.Printf("Archivos: %d, Directorios: %d, Omitidos: %d\n", filesCount, dirsCount, skippedCount)
	if dryRun {
		fmt.Println("(Modo dry-run: no se realizaron cambios)")
	}

	return nil
}
