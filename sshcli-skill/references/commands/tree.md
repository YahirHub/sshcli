# tree — Árbol de directorios

## Descripción
Muestra la estructura de directorios en forma de árbol.

## Sintaxis
```bash
sshcli tree [ruta] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-d, --depth` | Profundidad máxima |
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Árbol completo
sshcli tree "/var/www"

# Con profundidad
sshcli tree "/home" -d 2

# Mostrar solo directorios
sshcli tree "/var" -d 1
```

## Notas
- Útil para visualizar estructura de proyectos
- La profundidad `-d` limita la recursión
