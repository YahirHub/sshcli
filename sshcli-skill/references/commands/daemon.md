# daemon — Ejecutar en background

## Descripción
Ejecuta un comando como proceso en segundo plano.

## Sintaxis
```bash
sshcli daemon [comando] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-l, --log` | Archivo para guardar logs |
| `-n, --name` | Nombre identificador del daemon |
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Ejecutar en background
sshcli daemon "python /app/server.py"

# Con logs
sshcli daemon "node /app/index.js" --log "/var/log/app.log"

# Con nombre
sshcli daemon "npm start" --name myapp
```

## Notas
- El proceso persiste tras desconectarte
- Útil para servidores o scripts de larga duración
