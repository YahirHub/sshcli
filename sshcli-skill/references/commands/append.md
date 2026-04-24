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
# Agregar línea al final
sshcli append "/var/log/app.log" "INFO: Application started"

# Agregar contenido multilínea
sshcli append "/tmp/data.txt" "line1\nline2\nline3"
```

## Notas
- El archivo debe existir
- Usa `\n` para newlines
- No agrega espacio extra entre contenido existente y nuevo
