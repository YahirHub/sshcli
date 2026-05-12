# write — Crear o escribir archivos

## Descripción
Escribe contenido directamente a un archivo remoto. Crea directorios automáticamente si no existen.

## Sintaxis
```bash
sshcli write [ruta_remota] [contenido_opcional] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-x, --exec` | Hacer ejecutable (755) |
| `--chmod <perm>` | Permisos octales (default: 644) |
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Archivo simple
sshcli write "/app/file.txt" "hola"

# Multilínea con escapes
sshcli write "/app/file.txt" "linea1\nlinea2\n"

# Multilínea por stdin
printf 'a\nb\n' | sshcli write "/app/file.txt"

# Script ejecutable
sshcli write "/script/test.sh" "#!/bin/bash\necho OK\n" -x
```

## Notas
- El contenido puede pasarse como segundo argumento o por `stdin`.
- `\n`, `\t` y `\r` son interpretados correctamente.
- La flag `-x` otorga permisos 755.
- Crea directorios padre automáticamente.
