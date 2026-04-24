# git-status — Estado del repositorio

## Descripción
Muestra el estado del repositorio Git.

## Sintaxis
```bash
sshcli git-status [directorio] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
sshcli git-status "/home/user/project"
```

## Notas
- Muestra: rama actual, archivos modificados, staged, etc.
- Debe ejecutarse dentro de un repositorio Git
