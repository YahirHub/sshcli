# git-branch — Gestionar ramas

## Descripción
Gestiona ramas del repositorio Git.

## Sintaxis
```bash
sshcli git-branch [directorio] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-c, --create` | Crear nueva rama |
| `-d, --delete` | Eliminar rama |
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Listar ramas
sshcli git-branch "/repo"

# Crear rama
sshcli git-branch "/repo" -c "feature-new"

# Eliminar rama
sshcli git-branch "/repo" -d "feature-old"
```

## Notas
- Sin flags: lista ramas
- `*` = rama actual
