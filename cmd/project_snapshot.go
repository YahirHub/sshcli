package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var snapshotServer string

var snapshotCmd = &cobra.Command{
	Use:   "project-snapshot [ruta]",
	Short: "Genera un resumen completo del estado del proyecto",
	Long: `Recopila en un solo output:
  - Estructura de archivos (tree)
  - Últimos 5 commits de Git
  - Estado de servicios activos
  - Uso de recursos (disco/memoria)`,
	Args: cobra.ExactArgs(1),
	RunE: runSnapshot,
}

func init() {
	rootCmd.AddCommand(snapshotCmd)
	snapshotCmd.Flags().StringVarP(&snapshotServer, "server", "s", "", "Servidor específico")
}

func runSnapshot(cmd *cobra.Command, args []string) error {
	remotePath := cleanRemotePath(args[0])

	client, _, err := getClient(snapshotServer)
	if err != nil {
		return err
	}
	defer client.Close()

	fmt.Printf("======= SNAPSHOT DEL PROYECTO: %s =======\n\n", remotePath)

	// 1. Estructura de archivos
	fmt.Println("[1. ESTRUCTURA DE ARCHIVOS]")
	tree, _ := client.Run(fmt.Sprintf("tree -L 2 %s 2>/dev/null || ls -R %s | head -20", remotePath, remotePath))
	fmt.Println(tree)

	// 2. Últimos commits
	fmt.Println("\n[2. ÚLTIMOS 5 COMMITS]")
	commits, _ := client.Run(fmt.Sprintf("cd '%s' && git log -n 5 --oneline 2>/dev/null || echo 'No es un repo Git'", remotePath))
	fmt.Println(commits)

	// 3. Estado del Sistema
	fmt.Println("\n[3. RECURSOS DEL SISTEMA]")
	health, _ := client.Run("free -h && echo '' && df -h / | tail -1")
	fmt.Println(health)

	// 4. Servicios y Docker
	fmt.Println("\n[4. SERVICIOS Y CONTENEDORES]")
	docker, _ := client.Run("docker ps --format '{{.Names}} ({{.Status}})' 2>/dev/null || echo 'Docker no disponible'")
	fmt.Println(docker)

	fmt.Println("\n====================================================")
	return nil
}