# git-log — Historial de commits

## Descripción
Muestra el historial de commits.

## Sintaxis
```bash
sshcli git-log [directorio] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-n, --number` | Número de commits (default: 10) |
| `--oneline` | Formato compacto |
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Últimos 5 commits
sshcli git-log "/repo" -n 5

# Formato oneline
sshcli git-log "/repo" -n 5 --oneline
```

## Notas
- Default: muestra 10 commits
- `--oneline` = hash corto + mensaje
