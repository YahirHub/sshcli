# tail — Últimas líneas

## Descripción
Muestra las últimas líneas de un archivo.

## Sintaxis
```bash
sshcli tail [archivo] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-n, --lines` | Número de líneas (default: 20) |
| `-f, --follow` | Seguir el archivo (muestra cambios) |
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Últimas 50 líneas
sshcli tail "/var/log/syslog" -n 50

# Seguir log en tiempo real
sshcli tail "/var/log/app.log" -f
```

## Notas
- Útil para ver logs recientes
- La flag `-f` mantiene el archivo abierto
