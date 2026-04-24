# download — Descargar archivos

## Descripción
Descarga archivos del servidor remoto al local.

## Sintaxis
```bash
sshcli download [origen_remoto] [destino_local] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `--dry-run` | Simular sin descargar |
| `-e, --exclude` | Patrones a excluir |
| `--sync` | Modo sincronización |
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Descargar archivo
sshcli download "/tmp/file.txt" "/c/Users/Admin/downloads/"

# Excluir archivos
sshcli download "/var/www" "/backup/" --exclude "*.log"
```

## Notas
- Destino local usa formato MSYS2: `/c/Users/...`
- `--dry-run` muestra qué se descargaría sin hacerlo
- `--sync` elimina archivos locales que no existen remotamente
