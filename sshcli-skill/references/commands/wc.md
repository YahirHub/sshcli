# wc — Contar líneas, palabras y bytes

## Descripción
Cuenta líneas, palabras y bytes de un archivo.

## Sintaxis
```bash
sshcli wc [archivo] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-l, --lines` | Solo líneas |
| `-w, --words` | Solo palabras |
| `-c, --bytes` | Solo bytes |
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Contar todo
sshcli wc "/var/log/syslog"

# Solo líneas
sshcli wc "/app/data.csv" -l

# Solo bytes
sshcli wc "/var/log/syslog" -c
```

## Notas
- Output: `líneas palabras bytes archivo`
