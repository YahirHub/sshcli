# upload — Subir archivos

## Descripción
Sube un archivo o carpeta del local al servidor remoto.

## Sintaxis
```bash
sshcli upload [local] [remoto] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `--sync` | Modo sincronización |
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Subir archivo
sshcli upload "/c/Users/Admin/file.txt" "/tmp/file.txt"

# Subir directorio
sshcli upload "/c/Users/Admin/project" "/var/www/"

# Con sincronización
sshcli upload "/c/Users/Admin/project" "/var/www/" --sync
```

## Notas
- Ruta local usa formato MSYS2: `/c/Users/...`
- Normaliza rutas automáticamente
- `--sync` elimina archivos remotos que no existen localmente
