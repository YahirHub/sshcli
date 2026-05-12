# search-replace — Buscar y reemplazar

## Descripción
Busca y reemplaza texto dentro de un archivo en el servidor remoto.

## Sintaxis
```bash
sshcli search-replace [archivo] [buscar] [reemplazar] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-a, --all` | Reemplazar todas las ocurrencias |
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
sshcli search-replace "/etc/nginx.conf" "80" "8080"
sshcli search-replace "/app/config.py" "DEBUG=True" "DEBUG=False" -a
```

## Notas
- Por defecto reemplaza solo la primera ocurrencia.
- Usa `-a` para reemplazar todas.
- Operación segura: lee, modifica en memoria y sobreescribe.
