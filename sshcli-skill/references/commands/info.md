# info — Información de archivo

## Descripción
Muestra información detallada de un archivo o directorio.

## Sintaxis
```bash
sshcli info [ruta_remota] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
sshcli info "/var/log/syslog"
sshcli info "/home/user/file.txt"
```

## Notas
- Muestra: tamaño, permisos, dueño, grupo, fecha modificación
- Equivalente a `ls -la`
