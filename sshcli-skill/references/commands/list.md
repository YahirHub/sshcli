# list — Listar directorio

## Descripción
Lista archivos y directorios en una ruta remota.

## Sintaxis
```bash
sshcli list [ruta_remota] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Listar directorio
sshcli list "/var/www"

# Listar /tmp
sshcli list "/tmp"
```

## Notas
- Muestra archivos y subdirectorios
- No muestra archivos ocultos ( начина con `.`)
