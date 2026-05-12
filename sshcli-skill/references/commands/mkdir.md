# mkdir — Crear directorios

## Sintaxis
```bash
sshcli mkdir [ruta_remota] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-p, --parents` | Crear directorios padres si no existen |
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
sshcli mkdir "/home/user/newdir"
sshcli mkdir "/home/user/a/b/c" -p
```

## Notas
- `-p` crea directorios padre automáticamente.
