# git-add — Agregar al staging

## Descripción
Agrega archivos al área de staging.

## Sintaxis
```bash
sshcli git-add [directorio] [archivos...] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Agregar todos los archivos
sshcli git-add "/repo" "."

# Agregar archivo específico
sshcli git-add "/repo" "src/main.py"
```

## Notas
- Usa `.` para agregar todos los archivos
- Requiere al menos un archivo
