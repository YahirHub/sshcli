# kill — Terminar proceso

## Descripción
Termina un proceso en el servidor.

## Sintaxis
```bash
sshcli kill [pid_o_nombre] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `--signal` | Señal a enviar (default: 15/TERM) |
| `-s, --server` | Servidor específico |

## Señales Comunes
| Señal | Num | Descripción |
|-------|-----|-------------|
| TERM | 15 | Terminación graceful (default) |
| KILL | 9 | Terminación forzada |
| HUP | 1 | Reload configuración |
| INT | 2 | Interrupción (Ctrl+C) |

## Ejemplos
```bash
# Por PID
sshcli kill 12345

# Por nombre
sshcli kill nginx

# SIGKILL (forzado)
sshcli kill --signal 9 12345

# HUP (reload)
sshcli kill --signal HUP nginx
```

## Notas
- Por nombre usa `pkill`
- Default: señal 15 (TERM)
