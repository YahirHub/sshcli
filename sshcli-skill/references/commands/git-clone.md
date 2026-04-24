# git-clone — Clonar repositorio

## Descripción
Clona un repositorio Git.

## Sintaxis
```bash
sshcli git-clone [url_o_ruta] [destino] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Clonar desde URL
sshcli git-clone "https://github.com/user/repo.git" "/home/user/repo"

# Clonar repositorio local
sshcli git-clone "/old/repo" "/new/repo"
```

## Notas
- Puede clonar desde URL o ruta local
- Crea el directorio destino si no existe
