# copy — Copiar archivos

## Descripción
Copia archivos o directorios dentro del servidor remoto.

## Sintaxis
```bash
sshcli copy [origen_remoto] [destino_remoto] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Copiar archivo
sshcli copy "/etc/nginx.conf" "/etc/nginx.conf.bak"

# Copiar directorio completo
sshcli copy "/var/www/app" "/var/www/app.backup"
```

## Notas
- Ambas rutas son remotas
- Normaliza rutas automáticamente
- Sobrescribe si existe
