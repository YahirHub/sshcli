# wc — Contar líneas, palabras y bytes

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
sshcli wc "/var/log/syslog"
sshcli wc "/app/data.csv" -l
sshcli wc "/var/log/syslog" -c
```

## Notas
- Output típico: `líneas palabras bytes archivo`.
