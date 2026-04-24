# syntax-check — Validar sintaxis

## Descripción
Verifica la sintaxis de un archivo de código sin ejecutarlo.

## Sintaxis
```bash
sshcli syntax-check [archivo] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Verificar script Bash
sshcli syntax-check "/script/test.sh"

# Verificar Python
sshcli syntax-check "/app/main.py"

# Verificar Go
sshcli syntax-check "/app/main.go"
```

## Notas
- Soporta: `.sh`, `.py`, `.js`, `.go`, `.rb`, `.php`, `.pl`
- No soporta: `.txt`, `.json`, `.yaml` (extensiones no código)
