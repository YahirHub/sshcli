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
sshcli exec "ls -la /tmp"
sshcli exec -t "htop"
sshcli exec --no-tty "curl localhost:8080"
```

## Notas
- Modo normal: salida de texto.
- `-t` solicita un PTY remoto.
- Para una shell completa, usa `sshcli connect`.
