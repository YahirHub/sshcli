# grep — Buscar texto en archivos

## Descripción
Busca un patrón de texto dentro de archivos remotos.

## Sintaxis
```bash
sshcli grep [patron] [ruta_remota] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-i, --ignore-case` | Ignorar mayúsculas/minúsculas |
| `-r, --recursive` | Buscar recursivamente en directorios |
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Buscar patrón simple
sshcli grep "ERROR" "/var/log/syslog"

# Case insensitive
sshcli grep "error" "/var/log/app.log" -i

# Recursivo
sshcli grep "TODO" "/home/user/src" -r
```

## Notas
- Output: `linea:contenido`
- Útil para encontrar errores o configuraciones
