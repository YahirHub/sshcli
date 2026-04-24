# git-push — Push de cambios

## Descripción
Envía commits locales al repositorio remoto.

## Sintaxis
```bash
sshcli git-push [directorio] [rama] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Push al origin
sshcli git-push "/repo"

# Push de rama específica
sshcli git-push "/repo" "feature-branch"
```

## Notas
- Envía commits al remoto
- Puede requerir autenticación
