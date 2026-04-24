# write — Crear o escribir archivos

## Descripción
Escribe contenido directamente a un archivo remoto. Crea directorios automáticamente si no existen.

## Sintaxis
```bash
sshcli write [ruta_remota] [contenido] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-x, --exec` | Hacer ejecutable (755) |
| `--chmod <perm>` | Permisos octales (default: 644) |
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Crear script ejecutable
sshcli write "/script/test.sh" "#!/bin/bash\necho OK" -x

# Crear archivo con permisos específicos
sshcli write "/app/config.json" '{"key":"value"}' --chmod 600

# Crear directorios automáticamente
sshcli write "/deep/nested/path/file.txt" "content"
```

## Notas
- El contenido se pasa como segundo argumento
- Para contenido multilínea, usa `\n` para newlines
- La flag `-x` otorga permisos 755 (ejecutable)
- Crea directorios padre automáticamente
