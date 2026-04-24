# head — Primeras líneas

## Descripción
Muestra las primeras líneas de un archivo.

## Sintaxis
```bash
sshcli head [archivo] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-n, --lines` | Número de líneas (default: 20) |
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Primeras 10 líneas
sshcli head "/var/log/syslog" -n 10

# Default (20 líneas)
sshcli head "/etc/config.conf"
```

## Notas
- Default: 20 líneas
- Equivalente a `cat-lines archivo 1 N`
