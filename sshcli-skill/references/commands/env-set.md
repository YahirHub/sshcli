# env-set — Modificar archivo .env

## Descripción
Agrega o actualiza una variable en un archivo `.env`.

## Sintaxis
```bash
sshcli env-set [archivo_env] [KEY=value] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
sshcli env-set "/app/.env" "DATABASE_URL=postgres://localhost/db"
sshcli env-set "/app/.env" "DEBUG=false"
```

## Notas
- Si la variable existe, la actualiza.
- Si no existe, la agrega.
- Crea el archivo si no existe.
- Es recomendable verificar después con `sshcli read /app/.env`.
