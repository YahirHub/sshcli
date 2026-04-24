# delete-line — Eliminar líneas

## Descripción
Elimina una o más líneas de un archivo.

## Sintaxis
```bash
sshcli delete-line [archivo] [linea_inicio] [linea_fin] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Eliminar línea 15
sshcli delete-line "/app/main.py" 15 15

# Eliminar rango de líneas
sshcli delete-line "/app/main.py" 10 20
```

## Notas
- Especifica inicio y fin (puede ser el mismo número)
- Útil para eliminar código obsoleto o comentarios
