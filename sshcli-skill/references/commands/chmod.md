# chmod — Cambiar permisos

## Descripción
Cambia permisos de archivos o directorios.

## Sintaxis
```bash
sshcli chmod [permisos] [ruta_remota] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-r, --recursive` | Aplicar recursivamente |
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Hacer ejecutable
sshcli chmod +x "/script/test.sh"

# Permisos específicos
sshcli chmod 755 "/script/test.sh"
sshcli chmod 644 "/var/www/file.txt"

# Directorio recursivamente
sshcli chmod 755 "/var/www" -r
```

## Notas
- Formato: `+/-X` o octal (`755`, `644`, etc.)
- `+x` = agregar ejecutable
- `-x` = quitar ejecutable
