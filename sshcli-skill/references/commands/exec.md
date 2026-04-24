# exec — Ejecutar comandos

## Descripción
Ejecuta un comando en el servidor remoto.

## Sintaxis
```bash
sshcli exec [comando] [flags]
```

## Flags
| Flag | Descripción |
|------|-------------|
| `-t, --tty` | Habilitar modo interactivo (PTY) |
| `--no-tty` | Forzar modo normal |
| `-s, --server` | Servidor específico |

## Ejemplos
```bash
# Comando simple
sshcli exec "ls -la /tmp"

# Con interactive TTY
sshcli exec -t "htop"
sshcli exec -t "vim /etc/nginx.conf"

# Forzar modo no-TTY
sshcli exec --no-tty "curl localhost:8080"
```

## Notas
- Modo normal: salida de texto
- `-t` (TTY): para interactivos como htop, vim, apt
- Algunos comandos requieren rutas absolutas: `/bin/ls`
