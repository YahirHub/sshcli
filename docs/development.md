# Guía de Desarrollo - SSHCLI

## Requisitos

- Go 1.20 o superior
- Git

## Configuración del Entorno

```bash
git clone https://github.com/tu-usuario/sshcli.git
cd sshcli
go mod download
```

## Compilación

### Local
```bash
go build -o sshcli .
```

### Windows
```bash
go build -o sshcli.exe .
```

### Cross-compile para Linux desde Windows
```bash
set GOOS=linux
set GOARCH=amd64
go build -o sshcli .
```

### Cross-compile para Windows desde Linux
```bash
GOOS=windows GOARCH=amd64 go build -o sshcli.exe .
```

## Agregar un Nuevo Comando

### 1. Crear archivo en cmd/

```go
// cmd/mi_comando.go
package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var miComandoServer string

var miComandoCmd = &cobra.Command{
    Use:   "mi-comando [args]",
    Short: "Descripción corta en español",
    Long: `Descripción larga con ejemplos.

Ejemplos:
  sshcli mi-comando arg1
  sshcli mi-comando --server prod arg1

Casos de uso para agentes:
  - Caso 1
  - Caso 2`,
    Args: cobra.ExactArgs(1),
    RunE: runMiComando,
}

func init() {
    rootCmd.AddCommand(miComandoCmd)
    miComandoCmd.Flags().StringVarP(&miComandoServer, "server", "s", "", "Servidor específico a usar")
}

func runMiComando(cmd *cobra.Command, args []string) error {
    client, _, err := getClient(miComandoServer)
    if err != nil {
        return fmt.Errorf("error: %v", err)
    }
    defer client.Close()

    // Implementación
    output, err := client.Run("comando")
    if err != nil {
        return fmt.Errorf("error: %v", err)
    }

    fmt.Print(output)
    return nil
}
```

### 2. Convenciones

- Nombres de archivo: snake_case (`mi_comando.go`)
- Nombres de comando: kebab-case (`mi-comando`)
- Siempre incluir flag `-s, --server`
- Mensajes en español
- Incluir ejemplos en `Long`
- Documentar casos de uso para agentes

## Código Específico por OS

Para código que difiere entre Windows y Unix:

### Unix (connect_unix.go)
```go
//go:build !windows

package cmd

// Código específico Unix
```

### Windows (connect_windows.go)
```go
//go:build windows

package cmd

// Código específico Windows
```

## Testing

```bash
# Ejecutar tests
go test ./...

# Con cobertura
go test -cover ./...
```

## Estructura de Errores

```go
// Errores de configuración
return fmt.Errorf("configuración no encontrada: %v", err)

// Errores de conexión
return fmt.Errorf("error de conexión: %v", err)

// Errores de ejecución
return fmt.Errorf("error al ejecutar: %v", err)
```

## Funciones Comunes (cmd/common.go)

```go
// getClient obtiene cliente SSH para servidor especificado o activo
func getClient(serverName string) (*ssh.Client, *config.Server, error)
```

## Release

```bash
# Tag de versión
git tag v1.0.0
git push origin v1.0.0

# Build para múltiples plataformas
GOOS=linux GOARCH=amd64 go build -o dist/sshcli-linux-amd64 .
GOOS=darwin GOARCH=amd64 go build -o dist/sshcli-darwin-amd64 .
GOOS=windows GOARCH=amd64 go build -o dist/sshcli-windows-amd64.exe .
```

## Contribuir

1. Fork del repositorio
2. Crear rama: `git checkout -b feature/nueva-funcionalidad`
3. Commit: `git commit -m "feat: descripción"`
4. Push: `git push origin feature/nueva-funcionalidad`
5. Crear Pull Request

### Formato de Commits

```
tipo: descripción corta

Tipos:
- feat: nueva funcionalidad
- fix: corrección de bug
- docs: documentación
- refactor: refactorización
- test: tests
```
