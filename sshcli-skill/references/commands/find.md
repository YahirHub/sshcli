# find — Buscar archivos

## Descripción
Busca archivos y directorios usando patrones.

## Sintaxis
```bash
sshcli find [ruta_remota] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-n, --name` | Patrón de nombre a buscar |
| `-t, --type` | Tipo: `f`=archivo, `d`=directorio |
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Buscar por nombre
sshcli find "/var/log" -n "*.log"

# Solo archivos
sshcli find "/home" -n "*.py" -t f

# Solo directorios
sshcli find "/etc" -n "conf" -t d
```

## Notas
- Soporta wildcards: `*` = cualquier cosa, `?` = un caracter
- Patrones comunes: `*.log`, `*.conf`, `*.py`
