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
├── cmd/                    # Comandos CLI (58 archivos)
│   ├── root.go            # Comando raíz y configuración global
│   ├── common.go          # Funciones compartidas
│   ├── server.go          # Gestión de servidores
│   ├── config.go          # Configuración global
│   ├── connect.go         # Sesión SSH interactiva
│   ├── connect_unix.go    # Código específico Unix
│   ├── connect_windows.go # Código específico Windows
│   ├── exec.go            # Ejecución de comandos
│   ├── run.go             # Ejecutar scripts con intérprete automático
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
│   ├── find_code.go       # Buscar definiciones de código
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
│   ├── docker.go          # Gestión de contenedores Docker
│   ├── project_snapshot.go # Resumen completo del proyecto
│   └── status.go          # Estado de conexión
├── internal/               # Paquetes internos
│   ├── config/
│   │   └── config.go      # Gestión de configuración
│   ├── ssh/
│   │   └── client.go      # Cliente SSH/SFTP
│   └── paths/
│       └── paths.go        # Normalización de rutas (Unix/Windows)
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
- `GetActiveServer()` - Obtiene el servidor activo
- `AddServer()` / `RemoveServer()` - Gestión de servidores

### 3. Cliente SSH (internal/ssh/client.go)

Wrapper sobre las librerías SSH y SFTP:

```go
type Client struct {
    sshClient  *ssh.Client
    sftpClient *sftp.Client
}
```

Métodos principales:
- `Connect()` - Establece conexión SSH+SFTP
- `Run()` - Ejecuta comando remoto
- `WriteFile()` - Escribe archivo remoto
- `ReadFile()` - Lee archivo remoto
- `FileExists()` - Verifica existencia
- `IsDir()` - Verifica si es directorio

### 4. Normalización de Rutas (internal/paths/paths.go)

Manejo de rutas compatible con Git Bash/MSYS2 en Windows:

```go
// ToLocal - Convierte rutas Unix/MSYS2 a rutas nativas Windows
func ToLocal(p string) string

// ToRemote - Normaliza rutas para formato Linux absoluto
func ToRemote(p string) string
```

Esto permite que sshcli funcione correctamente cuando se ejecuta desde Git Bash en Windows.

### 5. Comandos (cmd/*.go)

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

### Comandos Avanzados

- `sshcli run /app/script.py` - Detecta automáticamente el intérprete (python3, node, bash, etc.)
- `sshcli project-snapshot /app` - Genera resumen completo: estructura, git, recursos, docker
- `sshcli find-code "func login" /app` - Busca definiciones de funciones/clases
- `sshcli docker ps` - Lista contenedores Docker activos

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
