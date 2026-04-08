package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sshcli",
	Short: "Cliente CLI SSH para agentes de IA",
	Long: `sshcli es un cliente CLI SSH diseñado para agentes de IA.
Permite ejecutar comandos remotos, subir/descargar archivos y carpetas
en múltiples servidores remotos de forma automatizada.

CONFIGURACIÓN INICIAL:
  sshcli server add mi-servidor --host 192.168.1.100 --user root --pass clave
  sshcli server add produccion --host prod.example.com --port 22 --user deploy --pass secreto

GESTIÓN DE SERVIDORES:
  sshcli server list              # Ver todos los servidores (* = activo)
  sshcli server use produccion    # Cambiar servidor activo
  sshcli server info mi-servidor  # Ver detalles de un servidor
  sshcli server remove mi-servidor

CONEXIÓN INTERACTIVA:
  sshcli connect                  # Abrir terminal SSH al servidor activo
  sshcli connect --server prod    # Conectar a servidor específico

EJECUCIÓN DE COMANDOS:
  sshcli exec "ls -la"            # Modo normal (salida de texto)
  sshcli exec -t htop             # Modo interactivo (pantalla completa)
  sshcli exec -t "apt install x"  # Con confirmaciones Y/N
  sshcli exec -t vim /etc/hosts   # Editores

HABILITAR -t POR DEFECTO:
  sshcli config set tty true      # Activar modo interactivo siempre
  sshcli config set tty false     # Desactivar
  sshcli exec --no-tty "ls"       # Forzar modo normal temporalmente

TRANSFERENCIA DE ARCHIVOS:
  sshcli upload archivo.txt /ruta/destino/
  sshcli download /ruta/remota/archivo.txt ./local/

EDICIÓN DE CÓDIGO:
  sshcli read /app/main.py                              # Leer archivo
  sshcli cat-lines /app/main.py 50 100                  # Leer líneas 50-100
  sshcli insert-line /app/main.py 1 "import os"         # Insertar línea
  sshcli search-replace /app/main.py "old" "new"        # Reemplazar texto

Todos los comandos soportan --server (-s) para especificar el servidor.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
