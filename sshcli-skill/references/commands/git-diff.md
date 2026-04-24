# git-diff — Ver cambios

## Descripción
Muestra los cambios en el repositorio.

## Sintaxis
```bash
sshcli git-diff [directorio] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `--staged` | Mostrar cambios staged |
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Cambios no staged
sshcli git-diff "/repo"

# Cambios staged
sshcli git-diff "/repo" --staged
```

## Notas
- Sin flags: muestra cambios no commitados
- `--staged`: muestra cambios en el área de staging
