# git-commit — Crear commit

## Descripción
Crea un commit con los cambios staged.

## Sintaxis
```bash
sshcli git-commit [directorio] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-a, --all` | Agregar todos los archivos modificados |
| `-m, --message` | Mensaje del commit (requerido) |
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Commit con mensaje
sshcli git-commit "/repo" -m "Add new feature"

# Commit con add automático
sshcli git-commit "/repo" -a -m "Fix bug"
```

## Notas
- `-m` es requerido
- `-a` agrega automáticamente archivos modificados (no nuevos)
