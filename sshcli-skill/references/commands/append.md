# append — Agregar contenido al final

## Descripción
Agrega texto al final de un archivo existente.

## Sintaxis
```bash
sshcli append [ruta_remota] [contenido] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
sshcli append "/var/log/app.log" "INFO: Application started"
sshcli append "/tmp/data.txt" "\nline2\nline3"
```

## Notas
- El archivo debe existir.
- `\n`, `\t` y `\r` se interpretan correctamente.
