package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sshcli",
	Short: "Cliente CLI SSH avanzado para agentes de IA",
	Long: `sshcli es un cliente SSH diseñado para que agentes de IA y modelos de lenguaje 
puedan administrar servidores Linux de forma atómica y profesional.

CONFIGURACIÓN INICIAL:
  sshcli server add mi-servidor --host 192.168.1.100 --user root --pass clave
  sshcli server list              # Ver servidores configurados

ANÁLISIS DE CONTEXTO (NUEVO):
  sshcli project-snapshot /app    # Resumen total: archivos, git, recursos y docker
  sshcli find-code "pattern" /app # Busca definiciones de funciones/clases
  sshcli tree /app -d 2           # Estructura de directorios

GESTIÓN DE CONTENEDORES (NUEVO):
  sshcli docker ps                # Listar contenedores activos
  sshcli docker logs[ID]         # Ver logs de un contenedor
  sshcli docker stats             # Ver consumo de CPU/RAM de contenedores
  sshcli docker exec [ID] "cmd"   # Ejecutar comando dentro de un contenedor

EJECUCIÓN Y CÓDIGO:
  sshcli exec "ls -la"            # Ejecución normal
  sshcli exec -t htop             # Modo interactivo (TTY)
  sshcli run /app/main.py         # Ejecutar script con intérprete automático
  sshcli syntax-check /app/file   # Validar sintaxis (Py, JS, Go, PHP, etc.)

EDICIÓN QUIRÚRGICA DE ARCHIVOS:
  sshcli read /app/file           # Leer contenido completo
  sshcli write /app/file          # Crear/sobreescribir desde stdin
  sshcli cat-lines /app/file 1 50 # Leer rango de líneas
  sshcli insert-line /app/file 5  # Insertar línea en posición X
  sshcli search-replace /app/file "old" "new" # Reemplazar texto

GESTIÓN DE GIT:
  sshcli git-status /app          # Estado del repositorio
  sshcli git-log /app             # Historial de cambios
  sshcli git-diff /app            # Ver diferencias actuales

SISTEMA Y RECURSOS:
  sshcli status                   # Verificar conexión al servidor
  sshcli memory                   # Uso de RAM
  sshcli disk /                   # Uso de disco
  sshcli ports -l                 # Puertos escuchando

Todos los comandos soportan el flag global --server (-s) para especificar el destino.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SilenceUsage = true
}