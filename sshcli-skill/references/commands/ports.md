# ports — Puertos de red

## Descripción
Muestra los puertos de red en uso en el servidor.

## Sintaxis
```bash
sshcli ports [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-l, --listen` | Solo puertos escuchando |
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Ver todos los puertos
sshcli ports

# Solo escuchando (default)
sshcli ports -l
```

## Notas
- Por defecto solo muestra puertos LISTEN
- Útil para verificar servicios activos
