# download — Descargar archivos

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
sshcli download "/tmp/file.txt" "/c/Users/Admin/downloads/file.txt"
sshcli download "/var/www" "/c/Users/Admin/backup" --exclude "*.log"
```

## Notas
- Destino local usa formato MSYS2: `/c/Users/...`.
