# server — Gestión de servidores

## Descripción
Gestiona los servidores SSH configurados.

## Sintaxis
```bash
sshcli server list         # Listar servidores
sshcli server add <nombre> --host <ip> --user <user> --pass <pass>
sshcli server use <nombre>  # Activar servidor
sshcli server info         # Info del servidor activo
sshcli server remove <nombre>  # Eliminar servidor
```

## Flags
| Flag | Descripción |
|------|-------------|
| `--host` | Dirección IP o hostname |
| `--user` | Usuario SSH |
| `--pass` | Contraseña |

## Ejemplos
```bash
# Ver servidores
sshcli server list

# Agregar servidor
sshcli server add "prod" --host "192.168.1.10" --user "root" --pass "secret"

# Usar servidor
sshcli server use "prod"

# Info del servidor activo
sshcli server info
```

## Notas
- `*` = servidor activo
- Requiere conexión SSH válida
