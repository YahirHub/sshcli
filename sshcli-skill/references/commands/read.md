# read — Leer contenido de archivo

## Descripción
Lee el contenido completo de un archivo remoto.

## Sintaxis
```bash
sshcli read [ruta_remota] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
sshcli read "/etc/config.conf"
sshcli read "/home/user/file.txt"
```

## Notas
- Para archivos grandes, usa `cat-lines` para especificar rango
- No soporta pipes ni redirecciones
