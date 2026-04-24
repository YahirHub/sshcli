# git-pull — Pull de cambios

## Descripción
Descarga y fusiona cambios del repositorio remoto.

## Sintaxis
```bash
sshcli git-pull [directorio] [rama] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Pull desde origin
sshcli git-pull "/repo"

# Pull de rama específica
sshcli git-pull "/repo" "feature-branch"
```

## Notas
- Fusiona cambios remotos en la rama actual
- Requiere conexión con repositorio remoto
