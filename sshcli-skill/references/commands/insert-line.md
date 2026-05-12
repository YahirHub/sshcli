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
# Insertar al inicio real del archivo
sshcli insert-line "/app/main.py" 1 "import os"

# Alias para inicio
sshcli insert-line "/app/main.py" 0 "import os"

# Insertar antes de la línea 10
sshcli insert-line "/app/config.py" 10 "DEBUG = True"
```

## Notas
- `1` inserta al inicio del archivo.
- `0` también se acepta como alias para inicio.
- `N` inserta antes de la línea `N`.
