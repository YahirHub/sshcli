# insert-line — Insertar línea

## Descripción
Inserta contenido en una línea específica. Las líneas siguientes se desplazan.

## Sintaxis
```bash
sshcli insert-line [archivo] [numero_linea] [contenido] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Insertar al inicio (línea 1)
sshcli insert-line "/app/main.py" 1 "import os"

# Insertar en línea 10
sshcli insert-line "/app/config.py" 10 "DEBUG = True"

# Insertar al final (línea 0)
sshcli insert-line "/app/main.py" 0 "# Nuevo comentario"
```

## Notas
- Línea 0 = inicio del archivo
- Las líneas existentes se desplazan hacia abajo
