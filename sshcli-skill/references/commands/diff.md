# diff — Comparar archivos

## Descripción
Compara un archivo local con uno remoto.

## Sintaxis
```bash
sshcli diff [archivo_local] [archivo_remoto] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-c, --context` | Líneas de contexto |
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Comparar archivos
sshcli diff "/c/local/file.py" "/remote/file.py"

# Con contexto
sshcli diff "/c/local/config.py" "/remote/config.py" -c 5
```

## Notas
- Formato local usa MSYS2: `/c/Users/...`
- Output en formato diff estándar
