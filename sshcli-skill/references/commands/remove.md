# remove — Eliminar archivos o directorios

## Descripción
Elimina archivos o directorios del servidor remoto.

## Sintaxis
```bash
sshcli remove [ruta_remota] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-f, --force` | Forzar eliminación sin confirmación |
| `-r, --recursive` | Eliminar directorios recursivamente |
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Eliminar archivo
sshcli remove "/tmp/test.txt"

# Eliminar directorio vacío
sshcli remove "/tmp/test-dir"

# Eliminar directorio con contenido
sshcli remove "/tmp/test-dir" -r

# Forzar eliminación
sshcli remove "/tmp/file.txt" -f
```

## Notas
- Directorios requieren `-r` para eliminarse con contenido
- No pide confirmación por defecto
- Usa `-f` para forzar sin confirmación
