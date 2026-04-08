# Arquitectura de SSHCLI

## Visión General

SSHCLI es un cliente SSH diseñado para agentes de IA, escrito en Go. La arquitectura sigue un patrón de comandos usando la librería Cobra.

## Estructura de Directorios

```
sshcli/
├── main.go                 # Punto de entrada
├── go.mod                  # Módulo Go
├── go.sum                  # Checksums de dependencias
├── README.md               # Documentación principal
├── cmd/                    # Comandos CLI
│   ├── root.go            # Comando raíz y configuración global
│   ├── common.go          # Funciones compartidas
│   ├── server.go          # Gestión de servidores
│   ├── config.go          # Configuración global
│   ├── connect.go         # Sesión SSH interactiva
│   ├── connect_unix.go    # Código específico Unix
│   ├── connect_windows.go # Código específico Windows
│   ├── exec.go            # Ejecución de comandos
│   ├── run.go             # Ejecutar scripts
│   ├── syntax_check.go    # Verificación de sintaxis
│   ├── read.go            # Leer archivos
│   ├── write.go           # Escribir archivos
│   ├── cat_lines.go       # Leer rango de líneas
│   ├── head.go            # Primeras líneas
│   ├── tail.go            # Últimas líneas
│   ├── insert_line.go     # Insertar línea
│   ├── delete_line.go     # Eliminar líneas
│   ├── replace_line.go    # Reemplazar línea
│   ├── search_replace.go  # Buscar y reemplazar
│   ├── append.go          # Agregar al final
│   ├── list.go            # Listar directorio
│   ├── tree.go            # Árbol de directorios
│   ├── mkdir.go           # Crear directorio
│   ├── copy.go            # Copiar archivos
│   ├── move.go            # Mover archivos
│   ├── remove.go          # Eliminar archivos
│   ├── chmod.go           # Cambiar permisos
│   ├── upload.go          # Subir archivos
│   ├── download.go        # Descargar archivos
│   ├── exists.go          # Verificar existencia
│   ├── info.go            # Información de archivo
│   ├── find.go            # Buscar archivos
│   ├── grep.go            # Buscar texto
│   ├── wc.go              # Contar líneas
│   ├── diff.go            # Comparar archivos
│   ├── git_*.go           # Comandos Git (10 archivos)
│   ├── ps.go              # Listar procesos
│   ├── kill.go            # Terminar procesos
│   ├── daemon.go          # Ejecutar en background
│   ├── service.go         # Gestionar servicios
│   ├── ports.go           # Ver puertos
│   ├── disk.go            # Uso de disco
│   ├── memory.go          # Uso de memoria
│   ├── env.go             # Variables de entorno
│   ├── env_set.go         # Establecer variable
│   └── status.go          # Estado de conexión
├── internal/               # Paquetes internos
│   ├── config/
│   │   └── config.go      # Gestión de configuración
│   └── ssh/
│       └── client.go      # Cliente SSH/SFTP
└── docs/                   # Documentación
    ├── architecture.md    # Este archivo
    ├── commands.md        # Referencia de comandos
    └── development.md     # Guía de desarrollo
```

## Componentes Principales

### 1. Main (main.go)

Punto de entrada que invoca el comando raíz:

```go
func main() {
    if err := cmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
```

### 2. Configuración (internal/config/config.go)

Gestiona la persistencia de servidores y configuración:

```go
type Config struct {
    ActiveServer string             // Servidor actualmente seleccionado
    Servers      map[string]*Server // Mapa de servidores
    DefaultTTY   bool               // TTY por defecto
}

type Server struct {
    Name     string
    Host     string
    Port     int
    User     string
    Password string
}
```

Funciones principales:
- `Load()` - Carga configuración desde ~/.sshcli.conf
- `Save()` - Guarda configuración
- `LoadOrCreate()` - Carga o crea nueva configuración

### 3. Cliente SSH (internal/ssh/client.go)

Wrapper sobre las librerías SSH y SFTP:

```go
type Client struct {
    sshClient  *ssh.Client
    sftpClient *sftp.Client
}
```

Métodos principales:
- `Connect()` - Establece conexión
- `Run()` - Ejecuta comando
- `WriteFile()` - Escribe archivo remoto
- `ReadFile()` - Lee archivo remoto

### 4. Comandos (cmd/*.go)

Cada comando sigue el patrón Cobra:

```go
var miCmd = &cobra.Command{
    Use:   "mi-comando [args]",
    Short: "Descripción corta",
    Long:  `Descripción larga con ejemplos...`,
    RunE:  runMiComando,
}

func init() {
    rootCmd.AddCommand(miCmd)
    miCmd.Flags().StringVarP(&variable, "flag", "f", "default", "descripción")
}

func runMiComando(cmd *cobra.Command, args []string) error {
    // Implementación
}
```

## Flujo de Ejecución

1. Usuario ejecuta: `sshcli exec "ls -la"`
2. `main.go` llama a `cmd.Execute()`
3. Cobra parsea argumentos y encuentra comando `exec`
4. `exec.go:runExec()` es invocado
5. Se carga configuración con `config.Load()`
6. Se obtiene servidor activo
7. Se crea conexión SSH con credenciales
8. Se ejecuta comando y se retorna salida

## Dependencias

```
github.com/spf13/cobra    # Framework CLI
github.com/pkg/sftp       # Cliente SFTP
golang.org/x/crypto/ssh   # Cliente SSH
golang.org/x/term         # Manejo de terminal
```

## Build Tags

Para compatibilidad multiplataforma:

- `connect_unix.go` - `//go:build !windows`
- `connect_windows.go` - `//go:build windows`

Esto permite manejar señales de terminal específicas de cada OS.
