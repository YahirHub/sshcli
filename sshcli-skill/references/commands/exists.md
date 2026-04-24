# exists — Verificar existencia

## Descripción
Verifica si un archivo o directorio existe.

## Sintaxis
```bash
sshcli exists [ruta_remota] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Verificar archivo
sshcli exists "/etc/nginx.conf"

# Verificar directorio
sshcli exists "/var/log"
```

## Notas
- Retorna `SI` si existe, `NO` si no existe
- Útil para verificar antes de operations críticas
