# git-checkout — Cambiar de rama

## Descripción
Cambia de rama o restaura archivos.

## Sintaxis
```bash
sshcli git-checkout [directorio] [rama_o_commit] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Cambiar a otra rama
sshcli git-checkout "/repo" "feature-branch"

# Volver a master
sshcli git-checkout "/repo" "master"
```

## Notas
- Cambia el HEAD a la rama especificada
- Requiere que la rama exista
