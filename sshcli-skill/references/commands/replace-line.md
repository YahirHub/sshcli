# replace-line — Reemplazar línea

## Descripción
Reemplaza completamente el contenido de una línea específica.

## Sintaxis
```bash
sshcli replace-line [archivo] [numero_linea] [nuevo_contenido] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Reemplazar línea 5
sshcli replace-line "/app/config.py" 5 "DEBUG = False"

# Reemplazar shebang
sshcli replace-line "/script/test.sh" 1 "#!/bin/bash"
```

## Notas
- La línea especificada se sustituye completamente
- No afecta otras líneas
