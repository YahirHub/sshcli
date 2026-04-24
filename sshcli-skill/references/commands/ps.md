# ps — Lista de procesos

## Descripción
Lista los procesos en ejecución en el servidor.

## Sintaxis
```bash
sshcli ps [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Lista completa
sshcli ps

# Filtrar con grep
sshcli ps | grep nginx
```

## Notas
- Equivalente a `ps aux`
- Muestra: USER, PID, %CPU, %MEM, VSZ, RSS, TTY, STAT, START, TIME, COMMAND
